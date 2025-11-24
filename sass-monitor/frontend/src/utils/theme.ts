import { ThemeMode, ThemeConfig } from '../types';

export class ThemeManager {
  private static readonly THEME_KEY = 'sass-monitor-theme';
  private static readonly DEFAULT_THEME: ThemeConfig = {
    mode: 'light',
    primary_color: '#1890ff',
    border_radius: 6,
  };

  /**
   * 获取当前主题配置
   */
  static getCurrentTheme(): ThemeConfig {
    try {
      const savedTheme = localStorage.getItem(this.THEME_KEY);
      return savedTheme ? JSON.parse(savedTheme) : this.DEFAULT_THEME;
    } catch (error) {
      console.error('Failed to parse theme config:', error);
      return this.DEFAULT_THEME;
    }
  }

  /**
   * 保存主题配置
   */
  static saveTheme(theme: ThemeConfig): void {
    try {
      localStorage.setItem(this.THEME_KEY, JSON.stringify(theme));
    } catch (error) {
      console.error('Failed to save theme config:', error);
    }
  }

  /**
   * 切换主题模式
   */
  static toggleThemeMode(): ThemeMode {
    const currentTheme = this.getCurrentTheme();
    const newMode: ThemeMode = currentTheme.mode === 'light' ? 'dark' : 'light';

    const newTheme: ThemeConfig = {
      ...currentTheme,
      mode: newMode,
    };

    this.saveTheme(newTheme);
    this.applyTheme(newTheme);

    return newMode;
  }

  /**
   * 设置主题模式
   */
  static setThemeMode(mode: ThemeMode): void {
    const currentTheme = this.getCurrentTheme();
    const newTheme = {
      ...currentTheme,
      mode,
    };

    this.saveTheme(newTheme);
    this.applyTheme(newTheme);
  }

  /**
   * 设置主色调
   */
  static setPrimaryColor(color: string): void {
    const currentTheme = this.getCurrentTheme();
    const newTheme = {
      ...currentTheme,
      primary_color: color,
    };

    this.saveTheme(newTheme);
    this.applyTheme(newTheme);
  }

  /**
   * 应用主题到DOM
   */
  static applyTheme(theme: ThemeConfig = this.getCurrentTheme()): void {
    const root = document.documentElement;

    // 设置主题模式
    root.setAttribute('data-theme', theme.mode);

    // 设置CSS变量
    this.setCSSVariables(theme);

    // 更新meta标签
    this.updateMetaTheme(theme.mode);
  }

  /**
   * 设置CSS变量
   */
  private static setCSSVariables(theme: ThemeConfig): void {
    const root = document.documentElement;

    // 主题颜色变量
    const colors = this.getThemeColors(theme.mode);

    Object.entries(colors).forEach(([key, value]) => {
      root.style.setProperty(`--theme-${key}`, value);
    });

    // 主色调
    root.style.setProperty('--theme-primary', theme.primary_color);

    // 圆角
    root.style.setProperty('--theme-border-radius', `${theme.border_radius}px`);

    // 浅蓝色主题变量
    if (theme.mode === 'light') {
      root.style.setProperty('--theme-primary', '#1890ff');
      root.style.setProperty('--theme-primary-light', '#40a9ff');
      root.style.setProperty('--theme-primary-dark', '#096dd9');
      root.style.setProperty('--theme-bg-primary', '#ffffff');
      root.style.setProperty('--theme-bg-secondary', '#fafafa');
      root.style.setProperty('--theme-bg-tertiary', '#f0f0f0');
      root.style.setProperty('--theme-text-primary', '#262626');
      root.style.setProperty('--theme-text-secondary', '#595959');
      root.style.setProperty('--theme-text-disabled', '#bfbfbf');
      root.style.setProperty('--theme-border-color', '#d9d9d9');
      root.style.setProperty('--theme-border-color-split', '#f0f0f0');
    } else {
      root.style.setProperty('--theme-primary', '#1890ff');
      root.style.setProperty('--theme-primary-light', '#40a9ff');
      root.style.setProperty('--theme-primary-dark', '#096dd9');
      root.style.setProperty('--theme-bg-primary', '#141414');
      root.style.setProperty('--theme-bg-secondary', '#1f1f1f');
      root.style.setProperty('--theme-bg-tertiary', '#262626');
      root.style.setProperty('--theme-text-primary', '#ffffff');
      root.style.setProperty('--theme-text-secondary', '#a6a6a6');
      root.style.setProperty('--theme-text-disabled', '#595959');
      root.style.setProperty('--theme-border-color', '#434343');
      root.style.setProperty('--theme-border-color-split', '#303030');
    }
  }

  /**
   * 获取主题颜色
   */
  private static getThemeColors(mode: ThemeMode): Record<string, string> {
    if (mode === 'light') {
      return {
        'bg-primary': '#ffffff',
        'bg-secondary': '#fafafa',
        'bg-tertiary': '#f0f0f0',
        'bg-hover': '#f5f5f5',
        'text-primary': '#262626',
        'text-secondary': '#595959',
        'text-disabled': '#bfbfbf',
        'border-primary': '#d9d9d9',
        'border-secondary': '#f0f0f0',
        'shadow': 'rgba(0, 0, 0, 0.1)',
        'card-shadow': '0 2px 8px rgba(0, 0, 0, 0.1)',
      };
    } else {
      return {
        'bg-primary': '#141414',
        'bg-secondary': '#1f1f1f',
        'bg-tertiary': '#262626',
        'bg-hover': '#303030',
        'text-primary': '#ffffff',
        'text-secondary': '#a6a6a6',
        'text-disabled': '#595959',
        'border-primary': '#434343',
        'border-secondary': '#303030',
        'shadow': 'rgba(0, 0, 0, 0.3)',
        'card-shadow': '0 2px 8px rgba(0, 0, 0, 0.3)',
      };
    }
  }

  /**
   * 更新meta主题色
   */
  private static updateMetaTheme(mode: ThemeMode): void {
    const themeColorMeta = document.querySelector('meta[name="theme-color"]');
    if (themeColorMeta) {
      const color = mode === 'light' ? '#ffffff' : '#141414';
      themeColorMeta.setAttribute('content', color);
    }
  }

  /**
   * 获取Ant Design主题配置
   */
  static getAntdTheme() {
    const theme = this.getCurrentTheme();

    return {
      token: {
        colorPrimary: theme.primary_color,
        borderRadius: theme.border_radius,
        wireframe: false,
      },
      algorithm: theme.mode === 'dark' ? [] : undefined,
      components: {
        Layout: {
          siderBg: theme.mode === 'dark' ? '#001529' : '#ffffff',
          triggerBg: theme.mode === 'dark' ? '#002140' : '#f0f0f0',
        },
        Menu: {
          darkItemBg: '#001529',
          darkSubMenuItemBg: '#000c17',
          darkItemSelectedBg: theme.primary_color,
        },
      },
    };
  }

  /**
   * 初始化主题
   */
  static initialize(): void {
    const theme = this.getCurrentTheme();
    this.applyTheme(theme);

    // 监听系统主题变化
    if (window.matchMedia) {
      const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
      mediaQuery.addListener(() => {
        // 用户可以根据系统主题自动切换
        // this.syncWithSystemTheme();
      });
    }
  }

  /**
   * 同步系统主题
   */
  static syncWithSystemTheme(): void {
    if (window.matchMedia) {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      this.setThemeMode(prefersDark ? 'dark' : 'light');
    }
  }

  /**
   * 获取主题模式的类名
   */
  static getThemeClassName(): string {
    const theme = this.getCurrentTheme();
    return `theme-${theme.mode}`;
  }

  /**
   * 重置为默认主题
   */
  static resetToDefault(): void {
    this.saveTheme(this.DEFAULT_THEME);
    this.applyTheme(this.DEFAULT_THEME);
  }

  /**
   * 预览主题
   */
  static previewTheme(theme: Partial<ThemeConfig>): void {
    const currentTheme = this.getCurrentTheme();
    const previewTheme = { ...currentTheme, ...theme };
    this.applyTheme(previewTheme);
  }

  /**
   * 取消主题预览
   */
  static cancelPreview(): void {
    const currentTheme = this.getCurrentTheme();
    this.applyTheme(currentTheme);
  }

  /**
   * 导出主题配置
   */
  static exportTheme(): string {
    const theme = this.getCurrentTheme();
    return JSON.stringify(theme, null, 2);
  }

  /**
   * 导入主题配置
   */
  static importTheme(themeJson: string): boolean {
    try {
      const theme = JSON.parse(themeJson) as ThemeConfig;

      // 验证主题配置
      if (this.validateTheme(theme)) {
        this.saveTheme(theme);
        this.applyTheme(theme);
        return true;
      }

      return false;
    } catch (error) {
      console.error('Failed to import theme:', error);
      return false;
    }
  }

  /**
   * 验证主题配置
   */
  private static validateTheme(theme: any): theme is ThemeConfig {
    return (
      typeof theme === 'object' &&
      theme !== null &&
      ['light', 'dark'].includes(theme.mode) &&
      typeof theme.primary_color === 'string' &&
      typeof theme.border_radius === 'number'
    );
  }

  /**
   * 获取预设主题
   */
  static getPresetThemes(): Array<{
    name: string;
    config: ThemeConfig;
    preview: string;
  }> {
    return [
      {
        name: '默认浅蓝',
        config: {
          mode: 'light',
          primary_color: '#1890ff',
          border_radius: 6,
        },
        preview: '#1890ff',
      },
      {
        name: '科技蓝',
        config: {
          mode: 'dark',
          primary_color: '#13c2c2',
          border_radius: 8,
        },
        preview: '#13c2c2',
      },
      {
        name: '活力橙',
        config: {
          mode: 'light',
          primary_color: '#fa8c16',
          border_radius: 6,
        },
        preview: '#fa8c16',
      },
      {
        name: '深空灰',
        config: {
          mode: 'dark',
          primary_color: '#722ed1',
          border_radius: 4,
        },
        preview: '#722ed1',
      },
    ];
  }
}

export default ThemeManager;