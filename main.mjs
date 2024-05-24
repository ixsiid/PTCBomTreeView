import { readFile } from 'fs/promises';

const [bom_file, sep, empty_line] = (args => {
	const file = args.filter(x => x.indexOf('-') != 0)[0];
	const sep = args.find(x => x === '-t') ? '\t' : ',';
	const empty_line = args.find(x => x === '-e');

	return [file, sep, empty_line];
})(process.argv.filter((_, i) => i >= 2));

console.error(`Start generate bom csv from ${bom_file}`);

readFile(bom_file)
	.then(f => f.toString())
	.then(text => {
		text = text.replaceAll('\r\n', '\n');
		text = text.replaceAll(/\n\s+[ABN]+\n/g, '\n');
		// データ行のスペーシング合わせ
		text = text.replaceAll('    サブアセンブリ', 'SubAssembly');
		text = text.replaceAll('  部品', 'Part');
		// タイトル行の英訳
		text = text.replaceAll('サブアセンブリ', 'SubAssembly');
		text = text.replaceAll('アセンブリ', 'Assembly');
		return text;
	})
	.then(text => {
		const sections = text.split('\n\n').map(section => {
			const s = section.split('\n');
			return {
				title: s.shift().split(' ').pop(),
				rows: s.map(l => l.split('*').map(x => x.trim())).map(([
					quantity,
					type,
					name,
					project,
					number,
					_class,
					supplier,
					material,
					finish,
				]) => ({
					quantity: parseInt(quantity),
					type,
					name,
					project,
					number,
					class: _class,
					supplier,
					material,
					finish,
				})),
			};
		});
		return sections;
	})
	.then(sections => {
		let s = sections.shift();
		const parts = sections.pop();

		const parse = x => {
			if (x.type === 'Part') {
				return parts.rows.find(r => r.name === x.name);
			} else if (x.type === 'SubAssembly') {
				if (x.material === 'PCBA') {
					return {
						type: 'PCBA',
						name: x.name,
						quantity: x.quantity,
					};
				} else {
					return {
						type: 'SubAssembly',
						name: x.name,
						quantity: x.quantity,
						rows: sections.find(s => s.title === x.name).rows.map(sx => parse(sx)),
					};
				}
			}
			console.error(x);
			throw 'Unknown type';
		};

		return {
			title: s.title,
			rows: s.rows.map(x => parse(x))
		};
	})
	.then(s => {
		console.log(['level', '', 'name', 'quantity', 'type', 'material', 'finish'].join(sep));
		// console.log('------------------------------------------')
		console.log(['', 'ASSY', 's.title'].join(sep));

		const print = (row, level) => {
			const n = '#'.repeat(level);
			if (row.type === 'Part') {
				console.log([n, 'PART', row.name, row.quantity, '', row.material, row.finish].join(sep));
			} else if (row.type === 'PCBA') {
				console.log([n, 'ASSY', row.name, row.quantity].join(sep));
			} else if (row.type === 'SubAssembly') {
				console.log([n, 'ASSY', row.name, row.quantity].join(sep));

				row.rows.forEach(r => print(r, level + 1));
				if (empty_line) console.log();
			} else {
				console.error(row);
				throw 'Unknwon type'
			}
		};

		s.rows.forEach(r => print(r, 1));
	});

