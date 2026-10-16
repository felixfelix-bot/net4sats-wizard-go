package main

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path"
	"time"

	"golang.org/x/crypto/ssh"
)

// sshConnect establishes an SSH session to the router.
//
// Auth chain, tried in order (a fresh-reset OpenWrt router ships root with
// an EMPTY password, while an already-configured one has the operator's
// password — the wizard must handle both):
//  1. Password(user-supplied)  — configured routers (v0.5.0 back-compat)
//  2. Password("")             — fresh routers, password auth
//  3. KeyboardInteractive      — fresh routers whose dropbear only accepts
//     (answers = password)       keyboard-interactive for the empty password
//  4. Default SSH keys          — key-provisioned routers, if present
func sshConnect(ip, password string) *ssh.Client {
	config := &ssh.ClientConfig{
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	auth := []ssh.AuthMethod{}
	if password != "" {
		auth = append(auth, ssh.Password(password))
	}
	auth = append(auth,
		ssh.Password(""),
		keyboardInteractiveAuth(password),
	)
	if signer := tryDefaultKeys(); signer != nil {
		auth = append(auth, ssh.PublicKeys(signer))
	}
	config.Auth = auth

	client, err := ssh.Dial("tcp", net.JoinHostPort(ip, "22"), config)
	if err != nil {
		return nil
	}
	return client
}

// keyboardInteractiveAuth answers every keyboard-interactive challenge with
// the given password ("" for a fresh router). The callback MUST return
// exactly one answer per question or x/crypto/ssh fails the auth attempt.
func keyboardInteractiveAuth(password string) ssh.AuthMethod {
	return ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
		answers := make([]string, len(questions))
		for i := range answers {
			answers[i] = password
		}
		return answers, nil
	})
}

// sshRun executes a command and returns combined output.
func sshRun(client *ssh.Client, cmd string) string {
	session, err := client.NewSession()
	if err != nil {
		return ""
	}
	defer session.Close()
	output, err := session.CombinedOutput(cmd)
	return string(output)
}

// sshUploadPipe writes binary data to the router via SSH stdin.
func sshUploadPipe(client *ssh.Client, data []byte, extractCmd string) string {
	session, err := client.NewSession()
	if err != nil {
		return ""
	}
	defer session.Close()
	session.Stdin = bytes.NewReader(data)
	output, err := session.CombinedOutput(extractCmd)
	return string(output)
}

// sshWriteFile writes content to a remote path via SSH (cat > path).
func sshWriteFile(client *ssh.Client, remotePath string, content []byte) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}

	if err := session.Start("cat > " + remotePath); err != nil {
		return err
	}

	_, err = stdin.Write(content)
	if err != nil {
		return err
	}
	stdin.Close()

	return session.Wait()
}

// sshDeployPortal writes the embedded portal files to the router.
func sshDeployPortal(client *ssh.Client, fs embed.FS, rootDir string) error {
	sshRun(client, "mkdir -p "+rootDir+"/assets "+rootDir+"/locales")

	entries, err := fs.ReadDir("portal")
	if err != nil {
		return fmt.Errorf("read portal embed: %w", err)
	}

	for _, entry := range entries {
		fullPath := "portal/" + entry.Name()
		if entry.IsDir() {
			subEntries, err := fs.ReadDir(fullPath)
			if err != nil {
				continue
			}
			sshRun(client, "mkdir -p "+path.Join(rootDir, entry.Name()))
			for _, sub := range subEntries {
				if sub.IsDir() {
					continue
				}
				data, err := fs.ReadFile(fullPath + "/" + sub.Name())
				if err != nil {
					continue
				}
				remotePath := path.Join(rootDir, entry.Name(), sub.Name())
				if err := sshWriteFile(client, remotePath, data); err != nil {
					return fmt.Errorf("write %s: %w", remotePath, err)
				}
			}
		} else {
			data, err := fs.ReadFile(fullPath)
			if err != nil {
				continue
			}
			remotePath := path.Join(rootDir, entry.Name())
			if err := sshWriteFile(client, remotePath, data); err != nil {
				return fmt.Errorf("write %s: %w", remotePath, err)
			}
		}
	}

	return nil
}

// fsDeployEntry is one planned operation in an embedded-FS deployment: the
// absolute remote path plus, for files, the bytes to write. Directories carry
// IsDir and no bytes (created with mkdir -p).
type fsDeployEntry struct {
	RemotePath string
	IsDir      bool
	Data       []byte
}

// planFSDeploy walks fsys under embedRoot and returns the ordered list of
// remote operations for remoteDir, honouring the embed-relative exclude
// list. It is pure (no SSH): this is the part of sshDeployFS that decides
// WHAT lands on the router, so it can be regression-tested without a router
// (see TestPortalDeployTargets).
//
// An embed-relative path in exclude is skipped entirely (file or directory).
// Embed paths always use "/", so remote paths are built with path.Join.
func planFSDeploy(fsys fs.FS, embedRoot string, remoteDir string, exclude ...string) ([]fsDeployEntry, error) {
	skip := make(map[string]bool, len(exclude))
	for _, e := range exclude {
		skip[e] = true
	}
	var plan []fsDeployEntry

	err := fs.WalkDir(fsys, embedRoot, func(fullPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip errors
		}
		// Compute relative path (embed paths always use /)
		relPath := fullPath
		prefix := embedRoot + "/"
		if len(fullPath) > len(prefix) && fullPath[:len(prefix)] == prefix {
			relPath = fullPath[len(prefix):]
		}
		if relPath == "" || relPath == embedRoot {
			return nil
		}
		if skip[relPath] {
			return nil
		}
		if d.IsDir() {
			plan = append(plan, fsDeployEntry{RemotePath: path.Join(remoteDir, relPath), IsDir: true})
			return nil
		}
		data, err := fs.ReadFile(fsys, fullPath)
		if err != nil {
			return nil
		}
		plan = append(plan, fsDeployEntry{RemotePath: path.Join(remoteDir, relPath), Data: data})
		return nil
	})
	if err != nil {
		return nil, err
	}
	return plan, nil
}

// sshDeployFS recursively deploys an embedded FS to a remote directory.
// Works with arbitrary nesting depth (e.g. assets/icon/colour/).
//
// exclude lists embed-relative paths (e.g. "index.html") that are skipped.
// Used when mirroring a bundle into a shared docroot so files owned by
// another service are left alone — e.g. OpenWrt's /www/index.html is the
// LuCI redirect page served on :8080 and must not be clobbered by a
// portal mirror landing in /www.
func sshDeployFS(client *ssh.Client, fsys embed.FS, embedRoot string, remoteDir string, exclude ...string) error {
	plan, err := planFSDeploy(fsys, embedRoot, remoteDir, exclude...)
	if err != nil {
		return err
	}

	sshRun(client, "mkdir -p "+remoteDir)

	for _, entry := range plan {
		if entry.IsDir {
			sshRun(client, "mkdir -p "+entry.RemotePath)
			continue
		}
		if err := sshWriteFile(client, entry.RemotePath, entry.Data); err != nil {
			return err
		}
	}

	return nil
}

// tryDefaultKeys attempts to load the default SSH key.
func tryDefaultKeys() ssh.Signer {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	for _, p := range []string{
		home + "/.ssh/id_ed25519",
		home + "/.ssh/id_rsa",
		home + "/.ssh/id_ecdsa",
	} {
		key, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err == nil {
			return signer
		}
	}
	return nil
}
