🐢 rm-buffer — Safe, Per-User rm Backups
=======================================

Because even the best engineers fat-finger `rm` sometimes.

---------------------------------------------------------------------------
💡 Why this project exists
---------------------------------------------------------------------------

Every Linux engineer has experienced that instant regret after typing
'rm -rf' too fast. rm-buffer exists to prevent that regret from becoming
a disaster.

It intercepts deletion only when you explicitly ask it to. When used
with -buffer, it quietly creates a compressed, restorable backup under
your home directory.

It’s not a replacement for rm — it’s a cautious companion.
Nothing changes unless you ask for it.

---------------------------------------------------------------------------
⚙️ How it works
---------------------------------------------------------------------------

When you run:
    rm -buffer myfile.txt

rm-buffer:
  1️⃣ Detects which user invoked it (even via sudo)
  2️⃣ Creates a timestamped tarball backup under ~/backupForRm/<basename>/
  3️⃣ Writes a small restore script (dependency.sh)
  4️⃣ Deletes the original only after the backup succeeds

Each backup lives in its own container — self-contained and restorable.

To restore:
    rm -extract <container-name>
or manually:
    tar -xzvf ~/backupForRm/<container>/<container>.tar.gz -C /target/path

---------------------------------------------------------------------------
📦 Project summary
---------------------------------------------------------------------------

Binary location:              /usr/local/bin/rm-buffer
Interactive shell helper:     /etc/profile.d/rm-buffer.sh
Admin helper (optional):      /usr/local/sbin/rm-buffer-enable
Per-user backup directory:    ~/backupForRm
Releases:                     https://github.com/itsadijmbt/rm-buffer-pkg/releases

---------------------------------------------------------------------------
🧰 Installation (Debian/Ubuntu)
---------------------------------------------------------------------------


1: curl -I "https://github.com/itsadijmbt/project-rm-buffer/releases/download/v1.0.2/rm-buffer_1.0.2_amd64.deb"

2: wget -O /tmp/rm-buffer.deb "https://github.com/itsadijmbt/rm-buffer-pkg/releases/download/v1.0.2/rm-buffer_1.0.2_amd64.deb"

3: wget -O /tmp/rm-buffer.deb.sha256 "https://github.com/itsadijmbt/rm-buffer-pkg/releases/download/v1.0.2/rm-buffer_1.0.2_amd64.deb.sha256"

4: sudo dpkg -i /tmp/rm-buffer.de

NOTE: if it fails open a new shell (profile.d scripts load automatically).

---------------------------------------------------------------------------
🚀 Usage overview
---------------------------------------------------------------------------

Backup-before-delete:
    rm -buffer <file>
    rm -b <file>             # shorthand

List available backups:
    rm -list
    rm -L

Extract a backup:
    rm -extract <container-name>
    rm -E <container-name>

---------------------------------------------------------------------------
🛠️ Admin: Optional global wrapper
---------------------------------------------------------------------------

To install a safe wrapper at /usr/local/bin/rm:
    sudo /usr/local/sbin/rm-buffer-enable

This script:
  - Backs up any existing /usr/local/bin/rm (e.g., rm.bak.<timestamp>)
  - Installs a wrapper that intercepts only -buffer/-list/-extract
  - For all other flags, calls /bin/rm directly

⚠️ Never overwrite /bin/rm or /usr/bin/rm.

---------------------------------------------------------------------------
🔁 Recovery example
---------------------------------------------------------------------------

$ rm -b /tmp/test.txt
Backed up: /home/user/backupForRm/test/test.tar.gz

$ rm -L
Available backups:
- test_2025-10-07T04:15:42Z

$ rm -E test_2025-10-07T04:15:42Z
Restored successfully.

---------------------------------------------------------------------------
🧩 Workflow integration
---------------------------------------------------------------------------

- Interactive-only: non-interactive scripts are unaffected.
- Each user has isolated backups at ~/backupForRm/.
- Supports custom cleanup or encryption via cron.
- Ideal for personal systems and dev environments.

---------------------------------------------------------------------------
🧑‍💻 Developer & packaging notes
---------------------------------------------------------------------------

Build the Go binary:
    cd src
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../bin/rm-buffer .

Create Debian package:
    mkdir -p pkg/usr/local/bin pkg/usr/local/sbin pkg/etc/profile.d pkg/DEBIAN
    cp bin/rm-buffer pkg/usr/local/bin/
    dpkg-deb --build pkg
    mv pkg.deb ../rm-buffer_1.0.2_amd64.deb
    sha256sum ../rm-buffer_1.0.2_amd64.deb > ../rm-buffer_1.0.2_amd64.deb.sha256

---------------------------------------------------------------------------
🩺 Troubleshooting
---------------------------------------------------------------------------

Problem: rm -list says "backup directory does not exist"
Cause:   No backups yet
Fix:     Run rm -b /path/to/file first

Problem: rm -buffer fails with "source does not exist"
Cause:   File already deleted
Fix:     Check path and retry

Problem: Normal user can’t see /root/backupForRm
Cause:   Backups are per-user
Fix:     Check ~/backupForRm instead

Remove everything:
    sudo dpkg -r rm-buffer
    rm -rf ~/backupForRm

---------------------------------------------------------------------------
🔒 Security & operational notes
---------------------------------------------------------------------------

- Default mode is safe; admin must opt-in for global wrapper.
- Backups are plain tarballs; treat as sensitive data.
- Add pruning, encryption, or remote archival as needed.
- HTTPS + sha256 verification recommended for all downloads.

---------------------------------------------------------------------------
🤝 Contributing
---------------------------------------------------------------------------

1️⃣ Fork the repo and create a feature branch.
2️⃣ Format code with `gofmt -w .`
3️⃣ Ensure `go build` passes.
4️⃣ Update pkg/DEBIAN/control for version bumps.
5️⃣ Submit a pull request with a short changelog.

---------------------------------------------------------------------------
🐢 Glumpy’s note
---------------------------------------------------------------------------

"Back up before you blow up."
— Glumpy the cautious turtle

---------------------------------------------------------------------------
📬 Contact / Support
---------------------------------------------------------------------------

GitHub issues and releases:
  https://github.com/itsadijmbt/rm-buffer-pkg


