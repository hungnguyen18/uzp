# Security Policy

## Overview

UZP-CLI is a security-focused tool designed to handle sensitive information such as API keys, passwords, and other secrets. We take security seriously and are committed to ensuring the protection of your data.

## Supported Versions

We provide security updates for the following versions:

| Version | Supported          | Go Requirement | Notes |
| ------- | ------------------ | -------------- | ----- |
| 1.0.x   | :white_check_mark: | Go 1.23.10+    | Current stable release |
| 0.9.x   | :white_check_mark: | Go 1.21+       | Legacy support until 2024-06-01 |
| < 0.9   | :x:                | -              | Please upgrade immediately |

> **⚠️ Security Notice**: Version 1.0.x requires Go 1.23.10+ to avoid [GO-2025-3750](https://pkg.go.dev/vuln/GO-2025-3750) affecting file operations.

**Recommendation:** Always use the latest stable version for the best security posture.

## Security Features

UZP-CLI implements multiple layers of security:

- **AES-256-GCM Encryption**: Industry-standard encryption for data at rest
- **scrypt Key Derivation**: Secure password-based key derivation (N=32768, r=8, p=1)
- **Secure File Permissions**: Vault files created with 0600 permissions (user-only access)
- **Memory Protection**: Sensitive data cleared from memory after use
- **No Password Storage**: Only password hashes stored for verification
- **Clipboard Security**: Automatic clipboard clearing with configurable TTL

## Reporting Security Vulnerabilities

### Quick Report
For urgent security issues, please email: **hungnguyen18.dev@gmail.com**

### Detailed Process

1. **DO NOT** create a public GitHub issue for security vulnerabilities
2. **DO** use GitHub's Private Security Advisory feature:
   - Go to: https://github.com/hungnguyen18/uzp-cli/security/advisories
   - Click "Report a vulnerability"
3. **Alternatively**, email: hungnguyen18.dev@gmail.com with:
   - Subject: `[SECURITY] UZP-CLI Vulnerability Report`
   - Detailed description of the vulnerability
   - Steps to reproduce
   - Your assessment of severity (Critical/High/Medium/Low)
   - Suggested mitigation if known

### What to Include

Please provide as much information as possible:

- **Version affected**: Which version(s) of UZP-CLI
- **Attack vector**: How the vulnerability can be exploited
- **Impact**: What data or functionality is at risk
- **Proof of concept**: Steps to reproduce (if safe to do so)
- **Environment**: OS, Node.js version, any relevant setup details

### Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial assessment**: Within 5 business days
- **Status updates**: Weekly until resolution
- **Fix timeline**: 
  - Critical: 24-72 hours
  - High: 1-2 weeks
  - Medium/Low: Next major release

## Disclosure Policy

We follow **coordinated disclosure** principles:

1. **Private reporting**: Issues reported privately first
2. **Investigation period**: Time to develop and test fixes
3. **User notification**: Security advisories published when fixes are available
4. **Recognition**: Security researchers credited (with permission)

### Timeline for Disclosure

- **Critical vulnerabilities**: 90 days maximum
- **High severity**: 120 days maximum  
- **Medium/Low severity**: 180 days maximum

## Security Best Practices for Users

### Installation Security

```bash
# ✅ Always verify package integrity
npm install -g uzp-cli

# ✅ Check package signature (if available)
npm audit

# ❌ Avoid installing from untrusted sources
```

### Usage Security

```bash
# ✅ Use strong master passwords
uzp init  # Use 12+ character passwords with mixed case, numbers, symbols

# ✅ Secure vault file permissions
chmod 600 ~/.uzp/uzp.vault

# ✅ Regular backups (encrypted)
cp ~/.uzp/uzp.vault ~/backup/uzp.vault.$(date +%Y%m%d)

# ❌ Never share your master password
# ❌ Don't store master password in scripts or files
```

### Environment Security

- **File Permissions**: Ensure `~/.uzp/` directory has 700 permissions
- **Backup Strategy**: Regularly backup your vault file securely
- **Network Security**: UZP-CLI works offline - no network access required
- **Multi-user Systems**: Each user should have their own vault
- **CI/CD**: Use environment variables, never commit vault files

## Known Security Considerations

### Current Limitations

1. **Memory dumps**: Sensitive data may briefly exist in memory
2. **Swap files**: Encrypted data might be written to swap (mitigate with encrypted swap)
3. **Process monitoring**: Admin users can inspect running processes
4. **Side-channel attacks**: Timing attacks theoretically possible during decryption

### Mitigations

- Use full-disk encryption on systems storing vault files
- Implement swap encryption or disable swap
- Run on systems with appropriate access controls
- Keep systems updated with latest security patches

## Security Hardening Guide

### For Individual Users

```bash
# Secure the UZP directory
chmod 700 ~/.uzp/
chmod 600 ~/.uzp/uzp.vault

# Set up encrypted backup
tar -czf - ~/.uzp/ | gpg --cipher-algo AES256 --compress-algo 1 \
  --symmetric --output uzp-backup-$(date +%Y%m%d).tar.gz.gpg

# Regular security audit
uzp list | wc -l  # Monitor number of stored secrets
```

### For Enterprise Users

- Deploy on hardened systems with minimal attack surface
- Implement centralized logging and monitoring
- Use configuration management for consistent security settings
- Consider hardware security modules (HSM) for additional protection
- Implement network segmentation where UZP-CLI is used

## Contact Information

- **Security Contact**: hungnguyen18.dev@gmail.com
- **General Issues**: https://github.com/hungnguyen18/uzp-cli/issues  
- **Security Advisories**: https://github.com/hungnguyen18/uzp-cli/security/advisories
- **GPG Key**: Available on request for encrypted communication

## Security Updates

Subscribe to security notifications:

- **GitHub Watch**: Enable security alerts for this repository
- **Release Notes**: Check CHANGELOG.md for security-related updates
- **Security Advisories**: GitHub will notify watchers of critical issues

---

**Last updated**: January 2025  
**Next review**: April 2025

Thank you for helping keep UZP-CLI and its users safe! 🔒
