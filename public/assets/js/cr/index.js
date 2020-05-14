class Element {
	constructor(elementId) {
		this.data = {
			element: document.getElementById(elementId),
			events: {},
			state: {},
			nextState: {},
			handleEvent: (from, e) => {
				e.preventDefault();
				const event = this.data.events[`${from}-${e.type}`];
				if (event)
					event.call({
						element: this.data.element,
						state: JSON.parse(JSON.stringify(this.data.state)),
						updateState: (key, value) => {
							this.data.nextState[key] = value;
						}
					}, e);

				const keys = Object.keys(this.data.nextState);
				if (keys.length > 0)
					for (const key of keys)
						this.data.state[key] = this.data.nextState[key];
				this.data.nextState = {};
			}
		}
	}

	on(event, func) {
		event = event.toLowerCase();

		this.data.events[`form-${event}`] = func;
		this.data.element[`on${event}`] = this.data.handleEvent.bind(this, 'form');
	}

	onInner(event, func) {
		event = event.toLowerCase();

		this.data.events[`inner-${event}`] = func;
		for (const child of this.data.element.children)
			child[`on${event}`] = this.data.handleEvent.bind(this, 'inner');
	}
}

class Form extends Element {
	constructor(elementId) {
		super(elementId);
	}
}

class ItemList extends Element {
	constructor(elementId, itemClass) {
		super(elementId);

		this.data.items = new Set();
		this.data.itemClass = itemClass;
		this.data.updateElement = async () => {
			while (this.data.element.children.length > 0)
				this.data.element.removeChild(this.data.element.firstChild);

			const items = Array.from(this.data.items);
			for (let i = 0; i < 20 && i < items.length; i++) {
				const item = new Item();
				await item.update(items[i], this.data.element);
			}
		}
	}

	async updateItems(newItems = []) {
		this.data.items = new Set();
		for (const item of newItems)
			this.data.items.add(item);
		await this.data.updateElement();
	}
}

class Item {
	constructor() {
		const h3 = document.createElement('h3');
		const img = document.createElement('img');
		const name = document.createElement('div');
		const desc = document.createElement('div');
		const input = document.createElement('input');
		const label = document.createElement('label');

		this.data = {
			state: {
				id: id => {input.value = id; input.setAttribute('id', id); label.setAttribute('for', id)},
				title: name => h3.textContent = name,
				imageLink: link => img.src = link,
				description: desc => desc.textContent = desc
			},
			element: label,
			input
		};

		h3.classList.add('title');
		name.appendChild(h3);

		desc.classList.add('desc');
		desc.appendChild(img);

		this.data.input.type = 'radio';
		this.data.input.name = 'anime';
		this.data.input.setAttribute('form', 'items');

		this.data.element.classList.add('item');
		this.data.element.appendChild(name);
		this.data.element.appendChild(desc);
	}

	async update(id, parent) {
		this.data.state.id(id);
		const itemData = await (fetch('/chewyroll/api/series/search/id/' + id).then(body => body.json()));

		if (itemData.name !== undefined) {
			this.data.state.title(itemData.name);
			this.data.state.imageLink(itemData.image);

			if (parent !== undefined) {
				parent.appendChild(this.data.input);
				parent.appendChild(this.data.element);
			}
			return true;
		}

		this.delete();
		return false;
	}

	delete() {
		if (this.data.element.parentElement) {
			this.data.element.parentElement.removeChild(this.data.element);
			this.data.element.parentElement.removeChild(this.data.input);
		}
	}
}

class ItemView extends Element {
	constructor(elementId, form) {
		super(elementId);

		this.data.lastUpdated = 0;
		this.data.radioForm = form;
		this.data.radioForm.on('change', async e => {
			await this.update(e.target.id);
		});
	}

	async update(id) {
		const time = Date.now();
		const title = this.data.element.getElementsByClassName('title').item(0);
		const headerImage = this.data.element.getElementsByClassName('header-image').item(0);
		const description = this.data.element.getElementsByClassName('description').item(0);
		const downloadProgress = this.data.element.getElementsByClassName('download-progress').item(0);

		const item = await (await fetch('/chewyroll/api/series/lookup/id/' + id)).json();
		if (this.data.lastUpdated < time) {
			title.textContent = item.title;
			headerImage.src = item.coverUri;
			description.textContent = item.description;
			downloadProgress.setAttribute('value', '0');
		}
	}
}

let data = {
	searchItemsUpdated: 0
};

window.onload = () => {
	const searchList = new ItemList('items');
	const searchForm = new Form('search');
	const itemView = new ItemView('info', searchList);

	const search = async value => {
		const time = Date.now();
		if (value === '' && data.searchItemsUpdated < time) {
			data.searchItemsUpdated = time;
			await searchList.updateItems([]);
		} else
			fetch('/chewyroll/api/series/search/name/' + value)
			.then(body => body.json())
			.then(json => {
				if (data.searchItemsUpdated < time) {
					data.searchItemsUpdated = time;
					return searchList.updateItems(json);
				}
			});
	}

	searchForm.on('submit', e => search(e.target[0].value));
	searchForm.onInner('input', e => search(e.target.value));
};

