import { store } from 'https://unpkg.com/hybrids@latest/src/index.js';
import wasm from '../api/wasm.js';
import ws from '../api/ws.js'
import ApplicationStore from './application-store.js';

const sleep = (milliseconds) => {
  return new Promise(resolve => setTimeout(resolve, milliseconds))
}

const AssetStore = {
    id: true,
    source: '',
    html: '',
    [store.connect]: {
        get: async (id) => {
            const app = store.get(ApplicationStore);
            const [assetType, name] = id.split(' ');
            const assetKey = { assetType: Number(assetType), name };
            let result;
            
            try {
                if (app.server) {
                    result = await ws.send('get', assetKey);
                } else {
                    result = await wasm.send('get', assetKey);
                }
            } catch(ex) {
                console.error('asset store get:', ex);
                return { id }
            }
            return {
                ...result,
                id: [result.id.assetType, result.id.name].join(' ')
            }
        },
        set: async (id, values, keys) => {id
            const app = store.get(ApplicationStore);
            const [assetType, name] = id.split(' ');
            const data = {
                ...values,
                id: {assetType: Number(assetType), name}
            };
            let result;

            if (app.server) {
                result = await ws.send('set', data);
            } else {
                result = await wasm.send('set', data);
                // TODO: Temporary workaround until 'hot reload' functionality built for WASM
                const main = await store.get(AssetStore, '0 main')
                store.clear(main)
            }
            return {
                ...result,
                id: [result.id.assetType, result.id.name].join(' ')
            }
        }
    }
}

export default AssetStore