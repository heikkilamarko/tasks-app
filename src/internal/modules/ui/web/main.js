import 'bootstrap/dist/css/bootstrap.css';
import './style.css';
import htmx from 'htmx.org';
import _hyperscript from 'hyperscript.org';
import { Modal, Toast, Tooltip } from 'bootstrap';
import { wsconnect } from '@nats-io/nats-core';

_hyperscript.browserInit();

window.app = Object.assign({}, window.app, {
	showConfirmModal,
	showToastMessage
});

document.addEventListener('DOMContentLoaded', async (_event) => {
	initBootstrap();

	if (!window.app.USER_ID) return;

	const nc = await wsconnect({
		servers: getWsUrl('/ws'),
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
	const data = msg.json();
	showToastMessage({
		type: 'warning',
		title: 'Task Expiring',
		text: data?.task?.name ?? '<no name>'
	});
}

function handleTaskExpiredMsg(msg) {
	const data = msg.json();
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

	const { title = '', text = '', details = '', type = 'info' } = config;

	if (!title && !text) return;

	const typeClasses = {
		info: {
			root: 'text-bg-primary',
			details: 'bg-primary-subtle text-primary-emphasis'
		},
		warning: {
			root: 'text-bg-warning',
			details: 'bg-warning-subtle text-warning-emphasis'
		},
		error: {
			root: 'text-bg-danger',
			details: 'bg-danger-subtle text-danger-emphasis'
		}
	};

	const classes = typeClasses[type] || typeClasses.info;

	const toastEl = document.createElement('div');
	toastEl.className = `toast fade border-0 ${classes.root}`;

	const toastBody = document.createElement('div');
	toastBody.className = 'toast-body d-flex flex-column';

	const titleEl = document.createElement('div');
	titleEl.className = 'fw-bold';
	titleEl.textContent = title;
	toastBody.appendChild(titleEl);

	const textEl = document.createElement('div');
	textEl.className = 'app-text-multiline';
	textEl.textContent = text;
	toastBody.appendChild(textEl);

	if (details) {
		const detailsEl = document.createElement('div');
		detailsEl.className = `app-text-multiline font-monospace mt-2 py-2 px-3 ${classes.details}`;
		detailsEl.textContent = details;
		toastBody.appendChild(detailsEl);
	}

	toastEl.appendChild(toastBody);
	toasterEl.appendChild(toastEl);

	const toast = new Toast(toastEl);

	toastEl.addEventListener(
		'hidden.bs.toast',
		() => {
			toast.dispose();
			toastEl.remove();
		},
		{ once: true }
	);

	toast.show();
}
