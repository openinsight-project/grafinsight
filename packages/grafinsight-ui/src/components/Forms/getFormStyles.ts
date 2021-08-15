import { stylesFactory } from '../../themes';
import { GrafInsightTheme } from '@grafinsight/data';
import { getLabelStyles } from './Label';
import { getLegendStyles } from './Legend';
import { getFieldValidationMessageStyles } from './FieldValidationMessage';
import { getButtonStyles, ButtonVariant } from '../Button';
import { ComponentSize } from '../../types/size';
import { getInputStyles } from '../Input/Input';
import { getCheckboxStyles } from './Checkbox';

export const getFormStyles = stylesFactory(
  (theme: GrafInsightTheme, options: { variant: ButtonVariant; size: ComponentSize; invalid: boolean }) => {
    return {
      label: getLabelStyles(theme),
      legend: getLegendStyles(theme),
      fieldValidationMessage: getFieldValidationMessageStyles(theme),
      button: getButtonStyles({
        theme,
        variant: options.variant,
        size: options.size,
      }),
      input: getInputStyles({ theme, invalid: options.invalid }),
      checkbox: getCheckboxStyles(theme),
    };
  }
);
