
import './App.css'
import { PrimeReactProvider } from 'primereact/api';
import "primereact/resources/themes/lara-light-cyan/theme.css";
import 'primeicons/primeicons.css';
import 'primeflex/primeflex.css';

import Layout from './components/Layout';

function App() {
    return (
        <PrimeReactProvider>
            <Layout />
        </PrimeReactProvider>
    )
}

export default App