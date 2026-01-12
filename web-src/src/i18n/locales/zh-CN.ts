export default {
  common: {
    save: '保存',
    saving: '保存中...',
    cancel: '取消',
    delete: '删除',
    add: '添加',
    edit: '编辑',
    confirm: '确认',
    apply: '应用',
    show: '显示',
    hide: '隐藏',
    confirmLogout: '确定要退出登录吗？',
    loading: '加载中...',
    seconds: '秒',
    days: '天',
    success: '成功',
    error: '错误',
    warning: '警告'
  },
  
  nav: {
    dashboard: '仪表盘',
    basic: '基础设置',
    integrations: '服务集成',
    mappings: '路径映射',
    ignore: '忽略配置',
    advanced: '高级设置',
    oauth: 'OAuth认证',
    logs: '实时日志',
    actions: '快捷操作',
    more: '更多'
  },
  
  header: {
    title: 'GD Watcher',
    subtitle: 'Google Drive 文件监控',
    save: '保存配置',
    saving: '保存中...',
    logout: '退出登录',
    language: '语言',
    theme: '主题色',
    customTheme: '自定义主题',
    appearance: '外观'
  },
  
  theme: {
    primaryColor: '主色调',
    secondaryColor: '次色调',
    opacity: '透明度'
  },
  
  dashboard: {
    systemStatus: '系统状态',
    online: '运行中',
    offline: '离线',
    uptime: '运行时间',
    activityStream: '活动流',
    historyCompletedTasks: '历史完成任务',
    todayCompletedTasks: '今日完成任务',
    cpuUsage: 'CPU 负载',
    memoryUsage: '内存负载'
  },
  
  sidebar: {
    collapse: '收起侧边栏',
    expand: '展开侧边栏'
  },
  
  login: {
    title: '欢迎回来',
    subtitle: '登录以继续使用 GD Watcher',
    username: '用户名',
    password: '密码',
    usernamePlaceholder: '请输入用户名',
    passwordPlaceholder: '请输入密码',
    submit: '登录',
    logging: '登录中...',
    error: '用户名或密码错误',
    footer: 'Google Drive 文件变更监控系统',
    refreshWallpaper: '更换壁纸'
  },
  
  logs: {
    title: '实时日志',
    autoScroll: '自动滚动',
    clear: '清空日志',
    empty: '暂无日志',
    newestFirst: '最新在前',
    oldestFirst: '最旧在前'
  },
  
  actions: {
    title: '快捷操作',
    testPath: '测试路径',
    testPathPlaceholder: '/云端硬盘/Test/1.mkv',
    test: '测试',
    sync: '触发同步',
    syncing: '同步中...',
    rcloneFull: '强刷 Rclone',
    rcloneFulling: '刷新中...',
    rebuildTree: '重建文件树',
    rebuildingTree: '重建中...',
    clearPanel: '清空面板',
    clearFiles: '清空日志文件',
    confirm: {
      rcloneFull: '确定要强制刷新所有 Rclone 实例吗？',
      rebuildTree: '确定要重建文件树吗？这可能需要一些时间。',
      clearFiles: '确定要删除所有日志文件吗？此操作不可恢复。'
    }
  },
  
  mappings: {
    regex: '正则表达式',
    replacement: '替换内容',
    regexPlaceholder: '^/云端硬盘/(.*)$',
    replacementPlaceholder: '/mnt/gd/$1',
    empty: '暂无映射规则',
    addFirst: '添加第一条规则'
  },
  
  cron: {
    placeholder: '0 0 * * *',
    invalid: '无效的 Cron 表达式',
    everyHour: '每小时',
    everyMinute: '每分钟',
    atTime: '{hour}:{minute}',
    custom: '自定义时间',
    daily: '每天',
    weekly: '每周',
    monthly: '每月',
    fields: {
      minute: '分',
      hour: '时',
      day: '日',
      month: '月',
      weekday: '周'
    }
  },
  
  panels: {
    basic: {
      title: '基础设置',
      description: '配置服务器端口、认证和 SSL 设置',
      server: '服务器配置',
      port: '监听端口',
      portHint: '服务监听端口，默认 8448',
      publicUrl: '公网地址',
      publicUrlPlaceholder: 'https://your-domain.com',
      publicUrlHint: '用于 Google Drive Webhook 回调',
      webhookPath: 'Webhook 路径',
      webhookPathHint: '用于 Google Drive Webhook 回调',
      learnMore: '了解更多',
      auth: '访问认证',
      username: '用户名',
      password: '密码',
      ssl: 'SSL 配置',
      enableSsl: '启用 HTTPS',
      sslCert: '证书路径',
      sslCertPlaceholder: '/path/to/cert.pem',
      sslKey: '私钥路径',
      sslKeyPlaceholder: '/path/to/key.pem'
    },
    
    integrations: {
      title: '服务集成',
      description: '配置 Google Drive、Rclone 和 Symedia 集成',
      googleDrive: 'Google Drive',
      qps: 'API QPS',
      qpsHint: '每秒 API 请求限制',
      learnMore: '了解更多',
      personalDriveName: '个人盘名称',
      personalDriveNamePlaceholder: '我的云端硬盘',
      rclone: 'Rclone 实例',
      rcloneHost: '主机地址',
      rcloneEndpoint: 'API 端点',
      waitForData: '等待数据同步',
      noRcloneInstances: '未配置 Rclone 实例',
      addRcloneInstance: '添加实例',
      symedia: 'Symedia 通知',
      symediaHost: '主机地址',
      symediaEndpoint: 'API 端点',
      symediaTemplate: '请求模板',
      symediaTemplateHint: '可用变量：',
      notifyUnmatched: '通知未匹配的路径',
      headers: '请求头',
      noHeaders: '未配置请求头'
    },
    
    mappings: {
      title: '路径映射',
      description: '配置 SA 和 Rclone 的路径转换规则',
      saMappings: 'SA 路径映射',
      saMappingsDesc: '将 Google Drive 路径转换为 Symedia 可识别的路径',
      rcloneMappings: 'Rclone 路径映射',
      rcloneMappingsDesc: '将 Google Drive 路径转换为 Rclone VFS 路径'
    },
    
    ignore: {
      title: '忽略配置',
      description: '配置需要忽略的 Google Drive 文件夹',
      placeholder: '输入 Google Drive 文件夹 ID',
      notePlaceholder: '添加备注（可选）',
      addNote: '点击添加备注',
      empty: '未配置忽略规则',
      emptyHint: '添加文件夹 ID 以忽略其中的文件变更'
    },
    
    advanced: {
      title: '高级设置',
      description: '配置同步参数和日志管理',
      syncSettings: '同步设置',
      debounce: '防抖延迟',
      debounceHint: '收到变更通知后等待的时间，避免频繁触发',
      logging: '日志配置',
      logDir: '日志目录',
      logLevel: '日志级别',
      logLevels: {
        quiet: '静默',
        info: '信息',
        debug: '调试'
      },
      enableLogSave: '保存日志到文件',
      logCleanup: '日志清理',
      enableLogCleanup: '启用自动清理',
      retentionDays: '保留天数',
      cleanupCron: '清理计划',
      cronLearnMore: 'Cron 表达式'
    },
    
    oauth: {
      title: 'OAuth 认证',
      description: '配置 Google OAuth 凭据',
      credentials: 'OAuth 凭据',
      clientId: 'Client ID',
      clientIdPlaceholder: '输入 Google OAuth Client ID',
      clientSecret: 'Client Secret',
      clientSecretPlaceholder: '输入 Client Secret',
      redirectUri: '重定向 URI',
      redirectUriPlaceholder: 'http://localhost:8448/api/oauth/callback',
      redirectUriHint: '需要在 Google Cloud Console 中配置',
      authorization: '授权操作',
      actionDesc: '配置好 OAuth 凭据后，点击下方按钮跳转到 Google 授权页面',
      goToGoogle: '前往 Google 授权',
      urlOpened: '已在新窗口打开授权页面',
      urlError: '获取授权链接失败，请检查 OAuth 配置',
      noClientId: '请先填写 Client ID'
    }
  },
  
  settings: {
    theme: '主题',
    themeLight: '浅色',
    themeDark: '深色',
    themeSystem: '跟随系统',
    language: '语言'
  },
  
  footer: {
    version: '版本',
    docs: '文档',
    github: 'GitHub'
  }
}
