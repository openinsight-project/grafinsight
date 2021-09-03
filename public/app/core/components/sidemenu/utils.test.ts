import { updateConfig } from '../../config';
import { getForcedLoginUrl } from './utils';

describe('getForcedLoginUrl', () => {
  it.each`
    appSubUrl          | url                    | expected
    ${''}              | ${'/whatever?a=1&b=2'} | ${'/whatever?a=1&b=2&forceLogin=true'}
    ${'/grafinsight'}      | ${'/whatever?a=1&b=2'} | ${'/grafinsight/whatever?a=1&b=2&forceLogin=true'}
    ${'/grafinsight/test'} | ${'/whatever?a=1&b=2'} | ${'/grafinsight/test/whatever?a=1&b=2&forceLogin=true'}
    ${'/grafinsight'}      | ${''}                  | ${'/grafinsight?forceLogin=true'}
    ${'/grafinsight'}      | ${'/whatever'}         | ${'/grafinsight/whatever?forceLogin=true'}
    ${'/grafinsight'}      | ${'/whatever/'}        | ${'/grafinsight/whatever/?forceLogin=true'}
  `(
    "when appUrl set to '$appUrl' and appSubUrl set to '$appSubUrl' then result should be '$expected'",
    ({ appSubUrl, url, expected }) => {
      updateConfig({
        appSubUrl,
      });

      const result = getForcedLoginUrl(url);

      expect(result).toBe(expected);
    }
  );
});
