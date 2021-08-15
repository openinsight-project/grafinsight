import React from 'react';
import { DataSourceApi, PanelData } from '@grafinsight/data';

interface InspectMetadataTabProps {
  data: PanelData;
  metadataDatasource?: DataSourceApi;
}
export const InspectMetadataTab: React.FC<InspectMetadataTabProps> = ({ data, metadataDatasource }) => {
  if (!metadataDatasource || !metadataDatasource.components?.MetadataInspector) {
    return <div>No Metadata Inspector</div>;
  }
  return <metadataDatasource.components.MetadataInspector datasource={metadataDatasource} data={data.series} />;
};
