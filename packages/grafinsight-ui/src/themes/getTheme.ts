import darkTheme from './dark';
import lightTheme from './light';
import { GrafInsightTheme } from '@grafinsight/data';

let themeMock: ((name?: string) => GrafInsightTheme) | null;

export const getTheme = (name?: string) =>
  (themeMock && themeMock(name)) || (name === 'light' ? lightTheme : darkTheme);

export const mockTheme = (mock: (name?: string) => GrafInsightTheme) => {
  themeMock = mock;
  return () => {
    themeMock = null;
  };
};
