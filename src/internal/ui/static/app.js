import { connect, StringCodec } from 'https://unpkg.com/nats.ws@1.17.0/esm/nats.js';

document.addEventListener('DOMContentLoaded', async (_event) => {
	const nc = await connect({
		servers: getWsUrl('/hub/v1'),
		token: 'S3c_r3t!'
	});
	const sc = StringCodec();
	const sub = nc.subscribe('tasks.ui.>');
	(async () => {
		for await (const msg of sub) {
			console.log(sc.decode(msg.data));
		}
	})();
});

export function getWsUrl(url) {
	return url?.startsWith('ws') ? url : `${location.origin.replace(/^http/, 'ws')}/${url.replace(/^\//, '')}`;
}
