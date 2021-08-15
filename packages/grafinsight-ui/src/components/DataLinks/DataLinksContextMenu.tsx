import React from 'react';
import { FieldConfig, LinkModel } from '@grafinsight/data';
import { selectors } from '@grafinsight/e2e-selectors/src';
import { css } from 'emotion';
import { WithContextMenu } from '../ContextMenu/WithContextMenu';
import { linkModelToContextMenuItems } from '../../utils/dataLinks';

interface DataLinksContextMenuProps {
  children: (props: DataLinksContextMenuApi) => JSX.Element;
  links: () => LinkModel[];
  config: FieldConfig;
}

export interface DataLinksContextMenuApi {
  openMenu?: React.MouseEventHandler<HTMLElement>;
  targetClassName?: string;
}

export const DataLinksContextMenu: React.FC<DataLinksContextMenuProps> = ({ children, links, config }) => {
  const linksCounter = config.links!.length;
  const getDataLinksContextMenuItems = () => {
    return [{ items: linkModelToContextMenuItems(links), label: 'Data links' }];
  };

  // Use this class name (exposed via render prop) to add context menu indicator to the click target of the visualization
  const targetClassName = css`
    cursor: context-menu;
  `;

  if (linksCounter > 1) {
    return (
      <WithContextMenu getContextMenuItems={getDataLinksContextMenuItems}>
        {({ openMenu }) => {
          return children({ openMenu, targetClassName });
        }}
      </WithContextMenu>
    );
  } else {
    const linkModel = links()[0];
    return (
      <a
        href={linkModel.href}
        onClick={linkModel.onClick}
        target={linkModel.target}
        title={linkModel.title}
        style={{ display: 'flex' }}
        aria-label={selectors.components.DataLinksContextMenu.singleLink}
      >
        {children({})}
      </a>
    );
  }
};
