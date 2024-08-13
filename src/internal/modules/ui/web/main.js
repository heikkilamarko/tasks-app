import 'bootstrap/dist/css/bootstrap.css';
import './style.css';
import htmx from 'htmx.org';
import _hyperscript from 'hyperscript.org';
import { Modal, Toast, Tooltip } from 'bootstrap';
import { wsconnect, JSONCodec } from '@nats-io/nats-core';

_hyperscript.browserInit();

window.app = Object.assign({}, window.app, {
	showConfirmModal,
	showToastMessage
});

const codec = JSONCodec();

document.addEventListener('DOMContentLoaded', async (_event) => {
	initBootstrap();

	if (!window.app.USER_ID) return;

	const nc = await wsconnect({
		servers: getWsUrl('/hub/v1'),
		name: 'ui',
		timeout: 120_000,
		maxReconnectAttempts: -1,
		waitOnFirstConnect: true
	});

	const sub = nc.subscribe(`task.${window.app.USER_ID}.>`);

	(async () => {
		for await (const msg of sub) {
			handleMsg(msg);
		}
	})();
});

document.body.addEventListener('htmx:responseError', (e) => {
	const {
		error,
		xhr: { responseText }
	} = e.detail;

	showToastMessage({
		type: 'error',
		title: 'ERROR',
		text: error || 'An error occurred.',
		details: responseText?.trim()
	});
});

document.body.addEventListener('htmx:sendError', (e) => {
	showToastMessage({
		type: 'error',
		title: 'ERROR',
		text: 'An error occurred. Reloading the page in 5 seconds...'
	});

	setTimeout(() => window.location.reload(), 5000);
});

htmx.onLoad((el) => {
	initBootstrap(el);
});

function initBootstrap(root) {
	const tooltips = (root ?? document).querySelectorAll('[data-bs-toggle="tooltip"]');
	[...tooltips].forEach((tooltip) =>
		Tooltip.getOrCreateInstance(tooltip, {
			customClass: 'app-tooltip',
			placement: 'left',
			delay: { show: 1000, hide: 100 }
		})
	);
}

function handleMsg(msg) {
	try {
		if (msg.subject.endsWith('.expiring')) {
			handleTaskExpiringMsg(msg);
		} else if (msg.subject.endsWith('.expired')) {
			handleTaskExpiredMsg(msg);
		} else {
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
		const modal = Modal.getOrCreateInstance(selector, { backdrop: 'static', keyboard: false });

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

	let classes = {
		root: 'text-bg-primary',
		details: 'bg-primary-subtle text-primary-emphasis'
	};

	switch (config.type) {
		case 'warning':
			classes = {
				root: 'text-bg-warning',
				details: 'bg-warning-subtle text-warning-emphasis'
			};
			break;
		case 'error':
			classes = {
				root: 'text-bg-danger',
				details: 'bg-danger-subtle text-danger-emphasis'
			};
			break;
	}

	const toastEl = document.createElement('div');
	toastEl.className = `toast fade border-0 ${classes.root}`;
	toastEl.innerHTML = `
		<div class="toast-body d-flex flex-column">
			<div class="fw-bold">${config.title}</div>
			<div class="app-text-multiline">${config.text}</div>
			${config.details ? `<div class="app-text-multiline font-monospace mt-2 py-2 px-3 ${classes.details}">${config.details}</div>` : ''}
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

	const toast = new Toast(toastEl);
	toast.show();
}
