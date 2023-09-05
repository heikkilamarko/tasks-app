import { connect, JSONCodec } from 'https://unpkg.com/nats.ws@1.17.0/esm/nats.js';

window.app = {
	showConfirmModal,
	showToastMessage
};

const codec = JSONCodec();

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
	const data = codec.decode(msg.data);
	showToastMessage({
		type: 'warning',
		title: 'Task Expiring',
		text: data?.task?.name ?? '<no name>'
	});
}

function handleTaskExpiredMsg(msg) {
	const data = codec.decode(msg.data);
	showToastMessage({
		type: 'error',
		title: 'Task Expired',
		text: data?.task?.name ?? '<no name>'
	});
}

function handleUnknownMsg(msg) {
	console.log('dropped unknown message', msg.subject);
}

function getWsUrl(url) {
	return url?.startsWith('ws') ? url : `${location.origin.replace(/^http/, 'ws')}/${url.replace(/^\//, '')}`;
}

function showConfirmModal(selector) {
	return new Promise((resolve) => {
		const modal = new bootstrap.Modal(selector, { backdrop: 'static', keyboard: false });

		modal._element.addEventListener(
			'confirmResult',
			(e) => {
				resolve(!!e.detail.answer);
				modal.hide();
			},
			{ once: true }
		);

		modal.show();
	});
}

function showToastMessage(config) {
	let typeClass = 'text-bg-primary';
	switch (config.type) {
		case 'warning':
			typeClass = 'text-bg-warning';
			break;
		case 'error':
			typeClass = 'text-bg-danger';
			break;
	}

	const toast = document.createElement('div');
	toast.className = `toast fade ${typeClass}`;
	toast.innerHTML = `
		<div class="toast-body">
			<strong class="me-auto">${config.title}</strong>:
			${config.text}
		</div>
	`;

	const toaster = document.getElementById('toaster');
	toaster.appendChild(toast);

	toast.addEventListener('hidden.bs.toast', () => toast.remove(), { once: true });

	new bootstrap.Toast(toast).show();
}
