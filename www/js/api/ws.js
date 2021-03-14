import { store } from 'https://unpkg.com/hybrids@latest/src/index.js';
import AssetStore from '../store/asset-store.js';

export default (function() {
    let id = 0;
    let socket;
    const promises = new Map();

    const send = (method, params) => {
        return new Promise((resolve, reject) => {
            const request = {
                jsonrpc: "2.0",
                method,
                params,
                id: ++id,
            }
            const attempt = () => {
                console.log(`Sending`, request)
                socket.send(JSON.stringify(request));
            }
            promises.set(request.id, {
                resolve,
                reject,
                interval: setInterval(attempt, 5000)
            })
            attempt()
        })
    };

    const socketMessageListener = (e) => {
        const response = JSON.parse(e.data)
        console.log(`Received`, response)
        
        if (response.jsonrpc === "2.0") {
            if (promises.has(response.id)) {
                const prom = promises.get(response.id)
                clearInterval(prom.interval)
                if (response.error) {
                    prom.reject(response.error)
                } else {
                    prom.resolve(response.result)
                }
                promises.delete(response.id)
            } else {
                console.error(response)
            }
        } else {
            const id = [response.id.assetType, response.id.name].join(' ') 
            const asset = store.get(AssetStore, id)
            store.clear(asset)
        }
    };
    
    const socketCloseListener = (e) => {
        if (socket) {
            console.info('Disconnected');
        }
    };

    const initialize = () => {
        const url = window.origin.replace("http", "ws") + '/ws';
        return new Promise((resolve, reject) => {
            socket = new WebSocket(url);
            socket.onopen = resolve;
            socket.onmessage = socketMessageListener;
            socket.onerror = (e) => {
                reject();
            }
            socket.onclose = socketCloseListener;
        });
    };

    return {
        initialize,
        send
    }
}())