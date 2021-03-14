
import { store } from 'https://unpkg.com/hybrids@latest/src/index.js';
import ws from '../api/ws.js'
import wasm from '../api/wasm.js'

const ApplicationStore = {
    server: false,
    ready: false,
};

window.addEventListener('load', async (event) => {
    try {
        await ws.initialize()
        store.set(ApplicationStore, {
            server: true,
            ready: true
        });
    } catch (ex) {
        await wasm.initialize()
        store.set(ApplicationStore, {
            ready: true
        });
    }
});

export default ApplicationStore