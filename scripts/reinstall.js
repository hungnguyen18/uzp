#!/usr/bin/env node

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

function printHeader() {
  console.log('🔄 UZP CLI Reinstaller');
  console.log('======================');
  console.log('');
}

function runCommand(command, description) {
  console.log(`🔧 ${description}...`);
  try {
    const output = execSync(command, { encoding: 'utf8', stdio: 'pipe' });
    if (output.trim()) {
      console.log(`   ${output.trim()}`);
    }
    console.log(`   ✅ Success`);
    return true;
  } catch (error) {
    console.log(`   ⚠️  ${error.message.split('\n')[0]}`);
    return false;
  }
}

function checkCurrentInstallation() {
  console.log('🔍 Checking current installation...');
  
  try {
    const version = execSync('uzp --version', { encoding: 'utf8', stdio: 'pipe' }).trim();
    console.log(`   Current version: ${version}`);
    return true;
  } catch (error) {
    console.log('   ℹ️  UZP not currently installed or not in PATH');
    return false;
  }
}

function getLatestVersion() {
  try {
    const npmInfo = execSync('npm view uzp-cli version', { encoding: 'utf8', stdio: 'pipe' }).trim();
    return npmInfo;
  } catch (error) {
    return 'latest';
  }
}

function reinstall() {
  printHeader();
  
  const hasCurrentInstall = checkCurrentInstallation();
  const latestVersion = getLatestVersion();
  
  console.log(`🎯 Target version: ${latestVersion}`);
  console.log('');
  
  // Step 1: Uninstall current version
  if (hasCurrentInstall) {
    console.log('📝 Step 1: Removing current installation');
    runCommand('npm uninstall -g uzp-cli', 'Uninstalling via npm');
    
    // Run our enhanced cleanup
    try {
      console.log('🧹 Running enhanced cleanup...');
      require('./uninstall.js');
    } catch (error) {
      console.log('   ⚠️  Enhanced cleanup failed, continuing...');
    }
  } else {
    console.log('📝 Step 1: No current installation found');
  }
  
  console.log('');
  
  // Step 2: Clear npm cache (helps with EEXIST issues)
  console.log('📝 Step 2: Clearing npm cache');
  runCommand('npm cache clean --force', 'Cleaning npm cache');
  
  console.log('');
  
  // Step 3: Fresh installation
  console.log('📝 Step 3: Installing fresh copy');
  const installSuccess = runCommand('npm install -g uzp-cli', 'Installing UZP CLI');
  
  console.log('');
  
  // Step 4: Verify installation
  console.log('📝 Step 4: Verifying installation');
  
  try {
    const newVersion = execSync('uzp --version', { encoding: 'utf8', stdio: 'pipe' }).trim();
    const uzpPath = execSync('which uzp', { encoding: 'utf8', stdio: 'pipe' }).trim();
    
    console.log(`   ✅ Version: ${newVersion}`);
    console.log(`   ✅ Location: ${uzpPath}`);
    
    // Test basic functionality
    try {
      execSync('uzp --help', { stdio: 'pipe' });
      console.log('   ✅ Basic functionality works');
    } catch (error) {
      console.log('   ⚠️  Basic functionality test failed');
    }
    
    console.log('');
    console.log('🎉 Reinstallation completed successfully!');
    console.log('💡 Try running: uzp --help');
    
  } catch (error) {
    console.log('   ❌ Verification failed');
    console.log('');
    console.log('🔧 Manual troubleshooting steps:');
    console.log('   1. Check your PATH includes npm global bin directory');
    console.log('   2. Try: npm config get prefix');
    console.log('   3. Ensure npm global directory is in your shell PATH');
    console.log('   4. Restart your terminal/shell');
    console.log('');
    console.log('🆘 If problems persist:');
    console.log('   - File an issue: https://github.com/hungnguyen18/uzp-cli/issues');
    console.log('   - Include your OS, Node.js version, and error messages');
  }
}

// Self-executing if run directly
if (require.main === module) {
  reinstall();
}

module.exports = { reinstall }; 