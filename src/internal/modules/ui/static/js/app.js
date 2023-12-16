import { connect, JSONCodec } from 'https://unpkg.com/nats.ws@1.20.0/esm/nats.js';

window.app = Object.assign({}, window.app, {
	showConfirmModal,
	showToastMessage
});

const codec = JSONCodec();

document.addEventListener('DOMContentLoaded', async (_event) => {
	initBootstrap();

	const nc = await connect({
		servers: getWsUrl(window.app.HUB_URL),
		user: 'ui',
		pass: 'S3c_r3t!',
		name: 'ui',
		timeout: 120_000,
		maxReconnectAttempts: -1,
		waitOnFirstConnect: true
	});

	const sub = nc.subscribe('tasks.ui.>');

	(async () => {
		for await (const msg of sub) {
			handleMsg(msg);
		}
	})();
});

document.body.addEventListener('htmx:responseError', (e) => {
	const {
		xhr: { status, statusText, responseText },
		error
	} = e.detail;

	showToastMessage({
		type: 'error',
		title: `${status} ${statusText}`,
		text: error,
		details: responseText
	});
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
		const modal = bootstrap.Modal.getOrCreateInstance(selector, { backdrop: 'static', keyboard: false });

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
	const toasterEl = document.getElementById('toaster');
	if (!toasterEl) return;

	let classes = 'border-primary text-primary';
	switch (config.type) {
		case 'warning':
			classes = 'border-warning text-warning';
			break;
		case 'error':
			classes = 'border-danger text-danger';
			break;
	}

	const toastEl = document.createElement('div');
	toastEl.className = `toast fade border ${classes}`;
	toastEl.innerHTML = `
		<div class="toast-body d-flex flex-column">
			<div class="fw-bold">${config.title}</div>
			<div class="app-text-multiline">${config.text}</div>
			${config.details ? `<div class="app-text-multiline"><code>${config.details}</code></div>` : ''}
		</div>
	`;

	toasterEl.appendChild(toastEl);

	toastEl.addEventListener(
		'hidden.bs.toast',
		() => {
			toast?.dispose();
			toastEl.remove();
		},
		{ once: true }
	);

	const toast = new bootstrap.Toast(toastEl);
	toast.show();
}
