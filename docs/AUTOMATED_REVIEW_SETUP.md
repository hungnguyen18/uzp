# Automated Review System Setup Guide

This guide is for **repository owners** to set up automated PR review systems for uzp-cli.

## 🚀 **Quick Setup Checklist**

```bash
# ✅ 1. GitHub Actions workflows are already committed
# ✅ 2. Labels configuration is ready
# ✅ 3. GitHub Copilot integration ready

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
3. Copilot will automatically start reviewing PRs when mentioned

### **Step 3: Configure Branch Protection**
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

### **Step 4: Test the System**
```bash
# Create a test PR to verify everything works:
git checkout -b test/automated-review-system
echo "// Test comment" >> cmd/root.go
git add . && git commit -m "test: verify automated review system"
git push origin test/automated-review-system

# Create PR and check:
# 1. ✅ Automated review comment appears
# 2. ✅ GitHub Copilot can be mentioned (@github-copilot review)
# 3. ✅ Labels are applied automatically
# 4. ✅ Security scanner runs automatically
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

### **Check Review System Health**

```bash
# 1. Monitor workflow runs
gh run list --workflow=pr-review.yml

# 2. Check recent PR comments  
gh pr list --state=all --limit=5

# 3. Verify system responses
# Look for comments from:
# - github-actions[bot] (automated review)
# - GitHub Copilot mentions working
# - Proper labels applied
```

### **Troubleshooting Common Issues**

| Issue | Solution |
|-------|----------|
| Copilot not responding | Check repo settings → Copilot enabled |
| Automated review not running | Check workflow permissions and triggers |
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

1. **✅ Create test PR** to verify automated system works
2. **📚 Announce to contributors** about GitHub Copilot integration
3. **🔔 Train contributors** on using `@github-copilot` commands
4. **📈 Monitor automated review effectiveness** and gather feedback
5. **🛠️ Iteratively improve** security patterns based on usage

**Congratulations!** Your repository now has a **streamlined automated review system** with GitHub Copilot that maintains security standards while keeping setup simple! 🚀

### **🎯 Benefits of This Simplified Approach:**

- **🟢 Less complexity** - Only GitHub Copilot + Automated Security Scanner
- **💰 Zero cost** - GitHub Copilot free for open source
- **⚡ Faster setup** - No external service integrations required
- **🔒 Same security** - Comprehensive automated security scanning
- **🤖 Smart AI** - GitHub Copilot provides excellent code review when needed 