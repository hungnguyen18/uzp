#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

const BINARY_PATH = path.join(__dirname, '..', 'bin', 'uzp');
const CACHE_DIR = path.join(os.homedir(), '.uzp-cache');

function getGlobalBinaryPaths() {
  const paths = [];
  
  try {
    // Get npm global bin directory
    const npmBin = execSync('npm bin -g', { encoding: 'utf8' }).trim();
    paths.push(path.join(npmBin, 'uzp'));
    
    // Common global paths
    const commonPaths = [
      '/usr/local/bin/uzp',
      path.join(os.homedir(), '.local/bin/uzp'),
      path.join(os.homedir(), '.npm-global/bin/uzp')
    ];
    
    // Add Node.js version-specific paths (for nvm users)
    if (process.env.NVM_DIR) {
      const nvmCurrent = process.env.NODE_VERSION || process.version;
      const nvmPath = path.join(process.env.NVM_DIR, 'versions', 'node', nvmCurrent, 'bin', 'uzp');
      commonPaths.push(nvmPath);
    }
    
    // Add paths from current Node.js installation
    if (process.execPath) {
      const nodeDir = path.dirname(process.execPath);
      commonPaths.push(path.join(nodeDir, 'uzp'));
    }
    
    paths.push(...commonPaths);
    
  } catch (error) {
    console.log('⚠️  Could not determine npm global directory, using common paths');
  }
  
  return [...new Set(paths)]; // Remove duplicates
}

function cleanupCache() {
  try {
    if (fs.existsSync(CACHE_DIR)) {
      fs.rmSync(CACHE_DIR, { recursive: true, force: true });
      console.log('🗑️  Cache directory cleaned');
    }
  } catch (error) {
    console.error('⚠️  Failed to cleanup cache:', error.message);
  }
}

function cleanupGlobalBinaries() {
  const globalPaths = getGlobalBinaryPaths();
  let cleanedCount = 0;
  
  console.log('🔍 Checking global binary locations...');
  
  for (const binPath of globalPaths) {
    try {
      if (fs.existsSync(binPath)) {
        console.log(`   Found: ${binPath}`);
        
        // Check if it's a symlink or regular file
        const stats = fs.lstatSync(binPath);
        if (stats.isSymbolicLink()) {
          fs.unlinkSync(binPath);
          console.log(`   ✅ Removed symlink: ${binPath}`);
        } else {
          fs.unlinkSync(binPath);
          console.log(`   ✅ Removed file: ${binPath}`);
        }
        cleanedCount++;
      }
    } catch (error) {
      console.log(`   ⚠️  Could not remove ${binPath}: ${error.message}`);
    }
  }
  
  if (cleanedCount === 0) {
    console.log('   ℹ️  No global binaries found to clean');
  } else {
    console.log(`🎉 Cleaned ${cleanedCount} global binary location(s)`);
  }
}

function uninstall() {
  console.log('🧹 Starting UZP cleanup...');
  
  try {
    // Clean up package binary
    if (fs.existsSync(BINARY_PATH)) {
      fs.unlinkSync(BINARY_PATH);
      console.log('✅ Package binary removed');
    }
    
    // Clean up global binaries
    cleanupGlobalBinaries();
    
    // Clean up cache
    cleanupCache();
    
    console.log('');
    console.log('✨ UZP cleanup completed!');
    console.log('💡 You can now reinstall with: npm install -g uzp-cli');
    
  } catch (error) {
    console.error('❌ Failed during cleanup:', error.message);
    console.log('');
    console.log('🔧 Manual cleanup suggestions:');
    console.log('   npm uninstall -g uzp-cli');
    console.log('   rm -f $(which uzp) 2>/dev/null');
    console.log('   npm install -g uzp-cli');
  }
}

uninstall(); 