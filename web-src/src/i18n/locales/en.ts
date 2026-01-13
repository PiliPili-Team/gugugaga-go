export default {
  common: {
    save: 'Save',
    saving: 'Saving...',
    cancel: 'Cancel',
    delete: 'Delete',
    add: 'Add',
    edit: 'Edit',
    confirm: 'Confirm',
    apply: 'Apply',
    show: 'Show',
    hide: 'Hide',
    confirmLogout: 'Are you sure you want to logout?',
    loading: 'Loading...',
    seconds: 'seconds',
    days: 'days',
    success: 'Success',
    error: 'Error',
    warning: 'Warning'
  },

  nav: {
    dashboard: 'Dashboard',
    basic: 'Basic Settings',
    integrations: 'Integrations',
    mappings: 'Path Mappings',
    ignore: 'Ignore Rules',
    advanced: 'Advanced',
    oauth: 'OAuth',
    logs: 'Live Logs',
    actions: 'Quick Actions',
    more: 'More'
  },

  header: {
    title: 'GD Watcher',
    subtitle: 'Google Drive File Monitor',
    save: 'Save Config',
    saving: 'Saving...',
    logout: 'Logout',
    language: 'Language',
    theme: 'Color Theme',
    customTheme: 'Custom Theme',
    appearance: 'Appearance'
  },

  theme: {
    primaryColor: 'Primary Color',
    secondaryColor: 'Secondary Color',
    opacity: 'Opacity'
  },

  dashboard: {
    systemStatus: 'System Status',
    online: 'Online',
    offline: 'Offline',
    uptime: 'Uptime',
    activityStream: 'Activity Stream',
    historyCompletedTasks: 'History Completed',
    todayCompletedTasks: 'Today Completed',
    cpuUsage: 'CPU Usage',
    memoryUsage: 'Memory Usage'
  },

  sidebar: {
    collapse: 'Collapse Sidebar',
    expand: 'Expand Sidebar'
  },

  login: {
    title: 'Welcome Back',
    subtitle: 'Sign in to continue to GD Watcher',
    username: 'Username',
    password: 'Password',
    usernamePlaceholder: 'Enter username',
    passwordPlaceholder: 'Enter password',
    submit: 'Sign In',
    logging: 'Signing in...',
    error: 'Invalid username or password',
    footer: 'Google Drive File Change Monitor',
    refreshWallpaper: 'Change Wallpaper'
  },

  logs: {
    title: 'Live Logs',
    autoScroll: 'Auto Scroll',
    clear: 'Clear Logs',
    empty: 'No logs yet',
    newestFirst: 'Newest First',
    oldestFirst: 'Oldest First'
  },

  actions: {
    title: 'Quick Actions',
    testPath: 'Test Path',
    testPathPlaceholder: '/My Drive/Test/1.mkv',
    test: 'Test',
    sync: 'Trigger Sync',
    syncing: 'Syncing...',
    rcloneFull: 'Full Rclone Refresh',
    rcloneFulling: 'Refreshing...',
    rebuildTree: 'Rebuild File Tree',
    rebuildingTree: 'Rebuilding...',
    clearPanel: 'Clear Panel',
    clearFiles: 'Clear Log Files',
    confirm: {
      rcloneFull: 'Are you sure you want to force refresh all Rclone instances?',
      rebuildTree: 'Are you sure you want to rebuild the file tree? This may take a while.',
      clearFiles: 'Are you sure you want to delete all log files? This cannot be undone.'
    }
  },

  mappings: {
    regex: 'Regex Pattern',
    replacement: 'Replacement',
    regexPlaceholder: '^/My Drive/(.*)$',
    replacementPlaceholder: '/mnt/gd/$1',
    empty: 'No mapping rules',
    addFirst: 'Add first rule'
  },

  cron: {
    placeholder: '0 0 * * *',
    invalid: 'Invalid Cron expression',
    everyHour: 'Every hour',
    everyMinute: 'Every minute',
    atTime: '{hour}:{minute}',
    custom: 'Custom time',
    daily: 'daily',
    weekly: 'weekly',
    monthly: 'monthly',
    fields: {
      minute: 'Min',
      hour: 'Hour',
      day: 'Day',
      month: 'Month',
      weekday: 'Week'
    }
  },

  panels: {
    basic: {
      title: 'Basic Settings',
      description: 'Configure server port, authentication, and SSL settings',
      server: 'Server Configuration',
      port: 'Listen Port',
      portHint: 'Service listening port, default 8448',
      publicUrl: 'Public URL',
      publicUrlPlaceholder: 'https://your-domain.com',
      publicUrlHint: 'Used for Google Drive Webhook callback',
      webhookPath: 'Webhook Path',
      webhookPathHint: 'Used for Google Drive Webhook callback',
      learnMore: 'Learn more',
      showPassword: 'Show password',
      hidePassword: 'Hide password',
      auth: 'Authentication',
      username: 'Username',
      password: 'Password',
      ssl: 'SSL Configuration',
      enableSsl: 'Enable HTTPS',
      sslCert: 'Certificate Path',
      sslCertPlaceholder: '/path/to/cert.pem',
      sslKey: 'Private Key Path',
      sslKeyPlaceholder: '/path/to/key.pem'
    },

    integrations: {
      title: 'Service Integrations',
      description: 'Configure Google Drive, Rclone, and Symedia integrations',
      googleDrive: 'Google Drive',
      qps: 'API QPS',
      qpsHint: 'API requests per second limit',
      learnMore: 'Learn more',
      personalDriveName: 'Personal Drive Name',
      personalDriveNamePlaceholder: 'My Drive',
      rclone: 'Rclone Instances',
      rcloneHost: 'Host Address',
      rcloneEndpoint: 'API Endpoint',
      waitForData: 'Wait for Data Sync',
      noRcloneInstances: 'No Rclone instances configured',
      addRcloneInstance: 'Add Instance',
      symedia: 'Symedia Notifications',
      symediaHost: 'Host Address',
      symediaEndpoint: 'API Endpoint',
      timeout: 'Timeout',
      timeoutHint: 'Max 120 seconds',
      symediaTemplate: 'Request Template',
      symediaTemplateHint: 'Available variables:',
      notifyUnmatched: 'Notify unmatched paths',
      headers: 'Headers',
      noHeaders: 'No headers configured'
    },

    mappings: {
      title: 'Path Mappings',
      description: 'Configure SA and Rclone path transformation rules',
      saMappings: 'SA Path Mappings',
      saMappingsDesc: 'Transform Google Drive paths to Symedia-recognizable paths',
      rcloneMappings: 'Rclone Path Mappings',
      rcloneMappingsDesc: 'Transform Google Drive paths to Rclone VFS paths'
    },

    ignore: {
      title: 'Ignore Rules',
      description: 'Configure Google Drive folders to ignore',
      placeholder: 'Enter Google Drive folder ID',
      notePlaceholder: 'Add note (optional)',
      addNote: 'Click to add note',
      empty: 'No ignore rules configured',
      emptyHint: 'Add folder IDs to ignore file changes within them'
    },

    advanced: {
      title: 'Advanced Settings',
      description: 'Configure sync parameters and log management',
      syncSettings: 'Sync Settings',
      debounce: 'Debounce Delay',
      debounceHint: 'Wait time after receiving change notifications to avoid frequent triggers',
      logging: 'Logging Configuration',
      logDir: 'Log Directory',
      logLevel: 'Log Level',
      logLevels: {
        quiet: 'Quiet',
        info: 'Info',
        debug: 'Debug'
      },
      enableLogSave: 'Save logs to file',
      logCleanup: 'Log Cleanup',
      enableLogCleanup: 'Enable auto cleanup',
      retentionDays: 'Retention Days',
      cleanupCron: 'Cleanup Schedule',
      cronLearnMore: 'Cron Expression'
    },

    oauth: {
      title: 'OAuth Authentication',
      description: 'Configure Google OAuth credentials',
      credentials: 'OAuth Credentials',
      clientId: 'Client ID',
      clientIdPlaceholder: 'Enter Google OAuth Client ID',
      clientSecret: 'Client Secret',
      clientSecretPlaceholder: 'Enter Client Secret',
      redirectUri: 'Redirect URI',
      redirectUriPlaceholder: 'http://localhost:8448/api/oauth/callback',
      redirectUriHint: 'Must be configured in Google Cloud Console',
      authorization: 'Authorization',
      actionDesc: 'After configuring OAuth credentials, click below to go to Google authorization page',
      goToGoogle: 'Go to Google Auth',
      urlOpened: 'Authorization page opened in new window',
      urlError: 'Failed to get authorization URL, please check OAuth configuration',
      noClientId: 'Please enter Client ID first'
    }
  },

  settings: {
    theme: 'Theme',
    themeLight: 'Light',
    themeDark: 'Dark',
    themeSystem: 'System',
    language: 'Language'
  },

  footer: {
    version: 'Version',
    docs: 'Docs',
    github: 'GitHub'
  }
}
