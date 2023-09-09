import { connect, JSONCodec } from 'https://unpkg.com/nats.ws@1.17.0/esm/nats.js';

window.app = {
	showConfirmModal,
	showToastMessage
};

const codec = JSONCodec();

document.addEventListener('DOMContentLoaded', async (_event) => {
	initBootstrap();

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

htmx.onLoad((el) => {
	initBootstrap(el);
});

function initBootstrap(root) {
	const tooltips = (root ?? document).querySelectorAll('[data-bs-toggle="tooltip"]');
	[...tooltips].forEach((tooltip) =>
		bootstrap.Tooltip.getOrCreateInstance(tooltip, {
			customClass: 'app-tooltip',
			placement: 'left',
			delay: { show: 1000, hide: 100 }
		})
	);
}

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
	const toaster = document.getElementById('toaster');
	if (!toaster) return;

	let typeClasses = 'border-primary text-primary';
	switch (config.type) {
		case 'warning':
			typeClasses = 'border-warning text-warning';
			break;
		case 'error':
			typeClasses = 'border-danger text-danger';
			break;
	}

	const toast = document.createElement('div');
	toast.className = `toast fade border ${typeClasses}`;
	toast.innerHTML = `
		<div class="toast-body d-flex flex-column">
			<div class="fw-bold">${config.title}</div>
			<div class="app-text-multiline">${config.text}</div>
		</div>
	`;

	toaster.appendChild(toast);

	toast.addEventListener('hidden.bs.toast', () => toast.remove(), { once: true });

	new bootstrap.Toast(toast).show();
}
