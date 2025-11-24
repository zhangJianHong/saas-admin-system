import { useState, useEffect } from 'react';
import { ThemeManager } from '../utils/theme';
import { ThemeMode, ThemeConfig } from '../types';

export const useTheme = () => {
  const [theme, setTheme] = useState<ThemeConfig>(ThemeManager.getCurrentTheme());
  const [isDark, setIsDark] = useState<boolean>(theme.mode === 'dark');

  useEffect(() => {
    // 初始化主题
    ThemeManager.initialize();

    // 监听主题变化
    const handleStorageChange = (e: StorageEvent) => {
      if (e.key === 'sass-monitor-theme') {
        const newTheme = ThemeManager.getCurrentTheme();
        setTheme(newTheme);
        setIsDark(newTheme.mode === 'dark');
      }
    };

    window.addEventListener('storage', handleStorageChange);
    return () => window.removeEventListener('storage', handleStorageChange);
  }, []);

  const toggleTheme = () => {
    const newMode = ThemeManager.toggleThemeMode();
    setIsDark(newMode === 'dark');
    setTheme(ThemeManager.getCurrentTheme());
    return newMode;
  };

  const setThemeMode = (mode: ThemeMode) => {
    ThemeManager.setThemeMode(mode);
    setIsDark(mode === 'dark');
    setTheme(ThemeManager.getCurrentTheme());
  };

  const setPrimaryColor = (color: string) => {
    ThemeManager.setPrimaryColor(color);
    setTheme(ThemeManager.getCurrentTheme());
  };

  const resetTheme = () => {
    ThemeManager.resetToDefault();
    setTheme(ThemeManager.getCurrentTheme());
    setIsDark(false);
  };

  return {
    theme,
    isDark,
    toggleTheme,
    setThemeMode,
    setPrimaryColor,
    resetTheme,
  };
};