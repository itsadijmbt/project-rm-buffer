rm-buffer — safe, per-user rm backups (plain text)
===============================================

A cautious, safe helper that only acts when you explicitly ask it to.
It archives files to a per-user backup location before deletion when
you use one of its explicit modes:
  -buffer  (or -b)
  -extract (or -E)
  -list    (or -L)

By default rm-buffer does NOT replace system `rm`. Admins may opt-in
to a safe wrapper; otherwise normal `rm` behavior is unchanged.

Mascot (Glumpy the cautious turtle):
   ____      
.-" +' "-.   
/  .-. .-. \  
/  /   Y   \ \
;  ;  .---.  ; ;
|  | (     ) | |
;  ;  '---'  ; ;
 \  \  '-'  / /
  '. `-._.-' .'  
    `-.___.-'










    
Glumpy the cautious turtle

-------------------------------------------------------------------------------
Quick summary
-------------------------------------------------------------------------------
- Binary installed:  /usr/local/bin/rm-buffer
- Interactive helper: /etc/profile.d/rm-buffer.sh (defines rm() for interactive shells only)
- Optional admin helper: /usr/local/sbin/rm-buffer-enable
- Backups location (per-user):  ~/backupForRm
- Releases: https://github.com/itsadijmbt/rm-buffer-pkg/releases

-------------------------------------------------------------------------------
Install (Debian/Ubuntu) — example (replace v1.0.2 with actual tag)
-------------------------------------------------------------------------------
wget -O /tmp/rm-buffer.deb "https://github.com/itsadijmbt/rm-buffer-pkg/releases/download/v1.0.2/rm-buffer_1.0.2_amd64.deb"
wget -O /tmp/rm-buffer.deb.sha256 "https://github.com/itsadijmbt/rm-buffer-pkg/releases/download/v1.0.2/rm-buffer_1.0.2_amd64.deb.sha256"
sha256sum -c /tmp/rm-buffer.deb.sha256
sudo dpkg -i /tmp/rm-buffer.deb

To enable the interactive helper in the current shell:
  source /etc/profile.d/rm-buffer.sh
Or open a new interactive terminal (profile.d is sourced at shell startup).

-------------------------------------------------------------------------------
Quick usage
-------------------------------------------------------------------------------
Note: plain rm behavior is preserved. Use these explicit modes to trigger backups.

Create a test file:
  printf 'hello\n' > /tmp/example.txt

Backup-before-delete:
  rm -buffer /tmp/example.txt
  rm -b /tmp/example.txt   # shorthand

List backups for current user:
  rm -list
  rm -L

Extract a backup (use a container name shown by rm -list):
  rm -extract <container-name>
  rm -E <container-name>

What happens:
- Backups are stored under:  ~/backupForRm/<basename>/
- Each container holds:  <basename>.tar.gz  and  dependency.sh
- dependency.sh is a small script that extracts the tarball back to the original location.

-------------------------------------------------------------------------------
Admin: opt-in global wrapper (recommended to be opt-in)
-------------------------------------------------------------------------------
To install a safe wrapper at /usr/local/bin/rm (optional):
  sudo /usr/local/sbin/rm-buffer-enable

This helper:
- Backs up any existing /usr/local/bin/rm (e.g. /usr/local/bin/rm.bak.<ts>)
- Installs a wrapper that only intercepts explicit modes (-buffer/-extract/-list)
- Otherwise delegates to /bin/rm

Do NOT overwrite /bin/rm or /usr/bin/rm.

-------------------------------------------------------------------------------
Recovering files
-------------------------------------------------------------------------------
1) Run:
   rm -list
   (find the backup container name)

2) Restore:
   rm -extract <container-name>
   OR manually:
   tar -xzvf ~/backupForRm/<container>/<container>.tar.gz -C /desired/location

-------------------------------------------------------------------------------
How it works (technical summary)
-------------------------------------------------------------------------------
- Interactive helper calls the binary only for explicit modes:
  /usr/local/bin/rm-buffer "{/pwd}" -buffer "/path/to/file"
- Binary normalizes paths, resolves invoking user's home (SUDO_USER used when appropriate),
  creates ~/backupForRm/<basename>/, tars the file or directory into <basename>.tar.gz,
  writes dependency.sh, moves artifacts into the backup folder, then deletes the original.
- -list reads the backup folder and prints a simple listing.
- -extract runs dependency.sh (or you can use tar manually).

-------------------------------------------------------------------------------
Developer / Maintainer: build & package (quick)
-------------------------------------------------------------------------------
# Build the Go binary:
cd src
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o ../bin/rm-buffer .

# Stage package layout:
cd ..
mkdir -p pkg/usr/local/bin pkg/usr/local/sbin pkg/etc/profile.d pkg/DEBIAN
cp -f ./bin/rm-buffer pkg/usr/local/bin/rm-buffer
chmod 0755 pkg/usr/local/bin/rm-buffer

# Ensure profile.d, wrapper, enable helper, and DEBIAN/* exist in pkg/

# Build deb:
dpkg-deb --build pkg
mv -f pkg.deb ../rm-buffer_1.0.2_amd64.deb
sha256sum ../rm-buffer_1.0.2_amd64.deb > ../rm-buffer_1.0.2_amd64.deb.sha256

-------------------------------------------------------------------------------
Troubleshooting & FAQ
-------------------------------------------------------------------------------
Q: rm -list says "backup directory does not exist"
A: No backups exist yet. Create one with: rm -buffer /path/to/file

Q: rm -buffer fails with "source does not exist"
A: The target file/directory no longer exists (maybe you already removed it).
   If you used -f and the file was missing, the command may skip silently.

Q: Normal user can't see /root/backupForRm
A: /root/backupForRm is root-owned and not visible to normal users. Per-user backups go to /home/<user>/backupForRm.

Q: How to remove package and backups?
A:
  sudo dpkg -r rm-buffer
  rm -rf ~/backupForRm

Q: Does this affect scripts or cronjobs?
A: No. The interactive helper is only sourced for interactive shells. Non-interactive processes remain unaffected.

-------------------------------------------------------------------------------
Security & operational notes
-------------------------------------------------------------------------------
- Default safe behavior: interactive-only; admin opt-in for global wrapper.
- Backups are plain tarballs in user home. Treat them as sensitive.
- Consider adding a prune/archive strategy for long-term operations.
- Prefer HTTPS release downloads and provide sha256 (and optionally GPG) signatures.

-------------------------------------------------------------------------------
Contributing
-------------------------------------------------------------------------------
1) Fork the repo and create a feature branch.
2) Make changes, run `gofmt -w .` and ensure `go build` passes.
3) Update pkg/DEBIAN/control version for packaging changes.
4) Open a pull request with a short changelog.

-------------------------------------------------------------------------------
License
-------------------------------------------------------------------------------
Add a LICENSE file (MIT recommended) in the repo root:
Example header:
  MIT License
  Copyright (c) 2025 Your Name

-------------------------------------------------------------------------------
Release & publishing notes
-------------------------------------------------------------------------------
- Tag releases (e.g. v1.0.2), upload rm-buffer_1.0.2_amd64.deb and rm-buffer_1.0.2_amd64.deb.sha256
  to GitHub Releases: https://github.com/itsadijmbt/rm-buffer-pkg/releases

- In README or release notes include:
  wget https://github.com/itsadijmbt/rm-buffer-pkg/releases/download/v1.0.2/rm-buffer_1.0.2_amd64.deb
  sha256sum -c rm-buffer_1.0.2_amd64.deb.sha256
  sudo dpkg -i rm-buffer_1.0.2_amd64.deb

-------------------------------------------------------------------------------
Contact / Support
-------------------------------------------------------------------------------
Open issues or PRs on GitHub:
  https://github.com/itsadijmbt/rm-buffer-pkg

-------------------------------------------------------------------------------
Final note from Glumpy:
-------------------------------------------------------------------------------
This tool is intentionally cautious — it does not change anything unless you ask for it.
Make backups, keep them tidy, and be careful. If you want, add pruning or encrypted remote
storage features as follow-up utilities.

End of README.txt

