import { store } from 'https://unpkg.com/hybrids@latest/src/index.js';
import AssetStore from '../store/asset-store.js';

export default (function(){
    const go = new Go();

    const send = async (method, params) => {
        let result
        switch(method) {
            case 'get':
                result = getAsset(params)
                if (!result.hasOwnProperty("id")) {
                    throw result;
                }
                return result
            case 'set':
                setAsset(params)
                result = getAsset(params.id)
                return result
            default:
                throw "Not implemented"
        }
    }

    const initialize = async () => {
        const result = await WebAssembly.instantiateStreaming(
            fetch(`www/lib.wasm`),
            go.importObject
        )
        go.run(result.instance)

        const data = await fetch("www/export.json")
        const assets = await data.json()

        await Promise.all(assets.map(asset => send('set', {id: asset.id, source: asset.source})))

        const asset = store.get(AssetStore, "0 main");
        store.clear(asset);
    }

    return {
        initialize,
        send
    }
}())