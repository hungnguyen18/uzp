#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const os = require('os');

const BINARY_PATH = path.join(__dirname, '..', 'bin', 'uzp');
const CACHE_DIR = path.join(os.homedir(), '.uzp-cache');

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

function uninstall() {
  try {
    if (fs.existsSync(BINARY_PATH)) {
      fs.unlinkSync(BINARY_PATH);
      console.log('✅ UZP binary removed successfully');
    }
    
    // Also cleanup cache
    cleanupCache();
    
  } catch (error) {
    console.error('❌ Failed to remove UZP binary:', error.message);
  }
}

uninstall(); 