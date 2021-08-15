import { dateTime, TimeRange } from '@grafinsight/data';
import { render } from '@testing-library/react';
import React from 'react';
import dark from '../../themes/dark';
import { UnthemedTimeRangePicker } from './TimeRangePicker';

const from = dateTime('2019-12-17T07:48:27.433Z');
const to = dateTime('2019-12-18T07:48:27.433Z');

const value: TimeRange = {
  from,
  to,
  raw: { from, to },
};

describe('TimePicker', () => {
  it('renders buttons correctly', () => {
    const container = render(
      <UnthemedTimeRangePicker
        onChangeTimeZone={() => {}}
        onChange={(value) => {}}
        value={value}
        onMoveBackward={() => {}}
        onMoveForward={() => {}}
        onZoom={() => {}}
        theme={dark}
      />
    );

    expect(container.queryByLabelText(/timepicker open button/i)).toBeInTheDocument();
  });
});
