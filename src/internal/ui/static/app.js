import { connect, StringCodec } from 'https://unpkg.com/nats.ws@1.17.0/esm/nats.js';

const sc = StringCodec();

document.addEventListener('DOMContentLoaded', async (_event) => {
	const nc = await connect({
		servers: getWsUrl('/hub/v1'),
		token: 'S3c_r3t!'
	});

	const sub = nc.subscribe('tasks.ui.>');

	(async () => {
		for await (const msg of sub) {
			handleMsg(msg);
		}
	})();
});

function handleMsg(msg) {
	try {
		switch (msg.subject) {
			case 'tasks.ui.expiring':
				handleTaskExpiringMsg(msg);
				break;
			case 'tasks.ui.expired':
				handleTaskExpiredMsg(msg);
				break;
			default:
				handleUnknownMsg(msg);
		}
	} catch (err) {
		console.error('message handling failed', err);
	}
}

function handleTaskExpiringMsg(msg) {
	const data = sc.decode(msg.data);
	Swal.fire({
		toast: true,
		position: 'top-end',
		icon: 'warning',
		title: 'Message',
		text: data,
		footer: '<a href="/ui">See expiring tasks</a>',
		showConfirmButton: false,
		timer: 5000
	});
}

function handleTaskExpiredMsg(msg) {
	const data = sc.decode(msg.data);
	Swal.fire({
		toast: true,
		position: 'top-end',
		icon: 'error',
		title: 'Message',
		text: data,
		footer: '<a href="/ui">See expired tasks</a>',
		showConfirmButton: false,
		timer: 5000
	});
}

function handleUnknownMsg(msg) {
	console.log('dropped unknown message', msg.subject);
}

function getWsUrl(url) {
	return url?.startsWith('ws') ? url : `${location.origin.replace(/^http/, 'ws')}/${url.replace(/^\//, '')}`;
}
