/** @type {import('tailwindcss').Config} */
export default {
	content: ['./src/**/*.{html,js,svelte,ts}'],
	theme: {
		extend: {
			fontFamily: {
				sans: [
					'-apple-system',
					'BlinkMacSystemFont',
					'"SF Pro Text"',
					'"SF Pro Display"',
					'"Segoe UI"',
					'Roboto',
					'Helvetica',
					'Arial',
					'sans-serif'
				],
			},
			colors: {
				ios: {
					blue: {
						light: '#0066CC', // Apple Web Link Blue
						dark: '#2997FF',  // Apple Web Dark Blue
					},
					gray: {
						100: '#F5F5F7', // Apple Web Light Gray Background
						200: '#E8E8ED',
						300: '#D2D2D7',
						400: '#86868B',
						500: '#1D1D1F', // Apple Web Dark Text
					},
					bg: {
						light: '#F5F5F7',
						dark: '#000000',
					},
					card: {
						light: '#FFFFFF',
						dark: '#1C1C1E',
					},
					separator: {
						light: 'rgba(60, 60, 67, 0.29)',
						dark: 'rgba(84, 84, 88, 0.65)'
					}
				}
			},
			borderRadius: {
				'ios-sm': '8px',
				'ios': '10px',
				'ios-md': '14px',
				'ios-lg': '20px',
			},
			boxShadow: {
				'ios': '0 4px 24px rgba(0, 0, 0, 0.04), 0 1px 2px rgba(0, 0, 0, 0.04)',
				'ios-modal': '0 20px 40px rgba(0, 0, 0, 0.1)',
			}
		}
	},
	plugins: []
};
