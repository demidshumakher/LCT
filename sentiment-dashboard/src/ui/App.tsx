import React from 'react';
import { Dashboard } from './Dashboard';

export const App: React.FC = () => {
	return (
		<div className="container">
			<div className="header">
				<h1 style={{ margin: 0 }}>Дашборд отзывов</h1>
			</div>
			<Dashboard />
			<div className="muted" style={{ marginTop: 12 }}>
				Данные взяты с{' '}
				<a href="https://www.banki.ru/" target="_blank" rel="noopener noreferrer">banki.ru</a>
			</div>
		</div>
	);
};



