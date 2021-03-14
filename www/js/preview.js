import { define, dispatch, html, store } from 'https://unpkg.com/hybrids@latest/src/index.js';
import AssetStore from './store/asset-store.js'

async function listStyles(host, e, attempt) {
  if (attempt === 3) {
    console.error('Unable to load styles')
  } else if (!attempt) {
    // TODO: Investigate more reliable means of detecting stylesheet load completion
    console.info('Retrying to load styles')
    setTimeout(() => {
      listStyles(host, e, (attempt || 0) + 1)
    }, 1000)
  } else {
    try {
      var allRules = [];
      const stylesheets = host.shadowRoot.styleSheets
      for (var sSheet = 0; sSheet < stylesheets.length; sSheet++) {
          var ruleList = stylesheets[sSheet].cssRules;
          for (var rule = 0; rule < ruleList.length; rule ++) {
            allRules.push( ruleList[rule].selectorText );
          }
      }
      const allRules2 = allRules.filter(a => a && a.match(/^\.[A-Za-z0-9\-]*$/)).map(a => a.slice(1))
      console.log(allRules2)
      dispatch(host, "classes", { detail: allRules2 })
    } catch(ex) {
      console.log('Unable to load styles, retrying...', ex)
      setTimeout(() => {
        listStyles(host, e, (attempt || 0) + 1)
      }, 1000)
    }
  }
}

const Preview = {
  assetType: 0,
  assetName: '',
  blueprint: false,
  asset: store(AssetStore, (host) => `${host.assetType} ${host.assetName}`),
  stylesheet: '',
  render: ({ asset, blueprint, stylesheet }) => html`
    <link href="${stylesheet}" crossorigin="anonymous" rel="stylesheet">
    ${listStyles}
    <style>
      .blueprint {
        background: radial-gradient(ellipse at center, #4c8ec4 0%, #0a3572 100%);
      }
      .blueprint * {
        background: rgba(0,0,0,0.1);
        color: white;
      }
      .blueprint *:hover {
        background: rgba(0,0,0,0.2);
      }
    </style>
    ${store.error(asset)}
    ${store.ready(asset) && html`
      <div id="preview" class="${{blueprint: blueprint}}" innerHTML="${asset.html}">
      </div>
    `}
  `,
};

define('x-preview', Preview);

export default Preview;