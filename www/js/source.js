import { define, dispatch, html, store } from 'https://unpkg.com/hybrids@latest/src/index.js';
import AssetStore from './store/asset-store.js'

const Source = {
  assetType: 0,
  assetName: '',
  asset: store(AssetStore, (host) => `${host.assetType} ${host.assetName}`),
  render: ({ asset }) => html`
    <style>
      :host {
        font-family: monospace;
      }
      #source {
        color: #60320b;
        white-space: pre-wrap;
      }
    </style>
    ${store.error(asset)}
    ${store.ready(asset) && html`
      <div id="source">${asset.html}</div>
    `}
  `,
};

define('x-source', Source);

export default Source;