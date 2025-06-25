# Automated Review System Setup Guide

This guide is for **repository owners** to set up automated PR review systems for uzp-cli.

## 🚀 **Quick Setup Checklist**

```bash
# ✅ 1. GitHub Actions workflows are already committed
# ✅ 2. CodeRabbit configuration is ready
# ✅ 3. Labels configuration is ready

# Run setup workflows:
```

### **Step 1: Setup Repository Labels**
```bash
# Trigger label setup workflow manually
gh workflow run setup-labels.yml
# Or push to main to trigger automatically
```

### **Step 2: Enable GitHub Copilot (Free for Open Source)**
1. Go to repository **Settings → Code security and analysis**
2. Enable **"GitHub Copilot"** for pull requests
3. Copilot will automatically start reviewing PRs

### **Step 3: Setup CodeRabbit AI**
1. Go to [coderabbit.ai](https://coderabbit.ai)
2. **Sign up with GitHub** (free for open source)
3. **Install CodeRabbit GitHub App** on your repository
4. CodeRabbit will automatically detect `.coderabbit.yaml` config

### **Step 4: Configure Branch Protection**
```bash
# Apply these settings via GitHub web interface:
# Settings → Branches → Add protection rule for main
```

**Branch Protection Settings:**
```yaml
Branch name pattern: main
✅ Require a pull request before merging
✅ Require status checks to pass before merging
  - ✅ ci / test
  - ✅ ci / build  
  - ✅ ci / lint
  - ✅ ci / security
✅ Require conversation resolution before merging
✅ Do not allow bypassing the above settings
❌ Allow force pushes: No
❌ Allow deletions: No
```

### **Step 5: Test the System**
```bash
# Create a test PR to verify everything works:
git checkout -b test/automated-review-system
echo "// Test comment" >> cmd/root.go
git add . && git commit -m "test: verify automated review system"
git push origin test/automated-review-system

# Create PR and check:
# 1. ✅ Automated review comment appears
# 2. ✅ CodeRabbit bot comments  
# 3. ✅ GitHub Copilot can be mentioned
# 4. ✅ Labels are applied automatically
```

---

## 🤖 **Individual Bot Setup**

### **GitHub Copilot Setup**

**Prerequisites:**
- Repository must be **public** (free for open source)
- Owner must have GitHub Pro/Team plan OR repository qualifies for Copilot for Open Source

**Setup:**
1. **Repository Settings** → **Code security and analysis**
2. **GitHub Copilot** → **Enable**
3. **Pull request reviews** → **Enable**

**Configuration:**
```yaml
# Already configured in .github/workflows/pr-review.yml
- name: Request Copilot Review
  uses: actions/github-script@v7
  with:
    script: |
      // Automatically requests Copilot review on PR open
```

**Usage:**
- Copilot automatically reviews PRs when requested by workflow
- Contributors can mention `@github-copilot review` for manual review
- Copilot provides security, performance, and quality feedback

---

### **CodeRabbit AI Setup**

**Prerequisites:**
- GitHub repository (public = free, private = paid)
- Admin access to repository

**Setup Steps:**
1. **Visit [coderabbit.ai](https://coderabbit.ai)**
2. **"Get Started"** → **"Sign in with GitHub"**
3. **Install CodeRabbit GitHub App**:
   - Select repositories: `uzp-cli` (or all repositories)
   - Grant permissions for PR review, comments, labels
4. **Configuration Auto-Detection**:
   - CodeRabbit will automatically detect `.coderabbit.yaml`
   - Custom rules for password manager security will be applied

**Verification:**
```bash
# Test CodeRabbit is working:
# 1. Create test PR
# 2. CodeRabbit should comment within 1-2 minutes
# 3. Look for @coderabbitai bot comments

# If not working, check:
# - GitHub App permissions
# - Repository access in CodeRabbit dashboard
# - .coderabbit.yaml syntax
```

**Available Commands:**
```bash
@coderabbitai help              # Show available commands
@coderabbitai review            # Request detailed review
@coderabbitai security          # Security-focused review
@coderabbitai performance       # Performance analysis  
@coderabbitai resolve           # Mark conversation resolved
```

---

### **Automated Security Scanner**

**Already Configured:**
- Runs automatically on every PR via GitHub Actions
- No additional setup required

**Features:**
- Hardcoded secret detection
- Weak cryptography detection
- Error handling validation
- Security gate for critical files

**Customization:**
```bash
# Edit security checks in:
.github/workflows/pr-review.yml

# Add custom security patterns:
- name: Security-Focused Review
  run: |
    # Add new patterns here
    if grep -r "NEW_SECURITY_PATTERN" --include="*.go" .; then
      security_issues="${security_issues}\n- ⚠️ Custom security issue"
    fi
```

---

## 🔧 **Advanced Configuration**

### **Custom Security Rules**

Edit `.coderabbit.yaml` to add password manager-specific rules:

```yaml
custom_rules:
  - name: "Check master password handling"
    pattern: "master.*password|password.*master"
    severity: high
    message: "Ensure master password is handled securely"
    
  - name: "Check memory clearing"
    pattern: "password|key|secret"
    severity: medium
    message: "Verify sensitive data is cleared from memory"
```

### **Workflow Customization**

Modify `.github/workflows/pr-review.yml` for your needs:

```yaml
# Change security check sensitivity
- name: Security-Focused Review
  run: |
    # Add more/fewer security patterns
    # Adjust severity levels
    # Custom notifications
```

### **Label Management**

Customize labels in `.github/workflows/setup-labels.yml`:

```yaml
const labels = [
  // Add custom labels for your workflow
  { name: 'custom-review', color: 'ff0000', description: 'Custom review needed' },
];
```

---

## 🔍 **Monitoring & Maintenance**

### **Check Review Bot Health**

```bash
# 1. Monitor workflow runs
gh run list --workflow=pr-review.yml

# 2. Check recent PR comments  
gh pr list --state=all --limit=5

# 3. Verify bot responses
# Look for comments from:
# - github-actions[bot]
# - coderabbit-ai[bot]  
# - GitHub Copilot mentions
```

### **Troubleshooting Common Issues**

| Issue | Solution |
|-------|----------|
| Copilot not responding | Check repo settings → Copilot enabled |
| CodeRabbit not commenting | Verify GitHub App permissions |
| Security scan false positives | Update regex patterns in workflow |
| Labels not applied | Run setup-labels workflow manually |

### **Performance Optimization**

```yaml
# Optimize workflow performance:
- name: Skip review for documentation-only PRs
  if: "!contains(github.event.pull_request.changed_files, '*.go')"
  
- name: Cache dependencies for faster runs
  uses: actions/cache@v3
  with:
    path: ~/go/pkg/mod
    key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
```

---

## 📊 **Success Metrics**

Track the effectiveness of your automated review system:

### **Weekly Review Health Check**
```bash
# Review metrics to track:
# 1. PR review response time (target: <2 minutes)
# 2. Security issues caught by automation (vs manual review)
# 3. Code quality improvements suggested
# 4. Contributor satisfaction with feedback quality
```

### **Monthly Optimization**
- Review false positive rates
- Update security patterns based on new threats
- Gather contributor feedback on bot usefulness
- Optimize workflow performance

---

## 🎯 **Next Steps**

After setup is complete:

1. **✅ Create test PR** to verify all systems work
2. **📚 Update contributor documentation** about new review process
3. **🔔 Announce to existing contributors** about automated review benefits
4. **📈 Monitor effectiveness** and gather feedback
5. **🛠️ Iteratively improve** based on usage patterns

**Congratulations!** Your repository now has a comprehensive automated review system that maintains security standards while scaling contributor onboarding! 🚀 