import './App.css'
import { PrimeReactProvider } from 'primereact/api';


import 'primeicons/primeicons.css';

import "primereact/resources/themes/lara-light-cyan/theme.css";
//import 'primereact/resources/themes/mdc-dark-indigo/theme.css';
//import 'primereact/resources/themes/bootstrap4-dark-purple/theme.css'

import LeftMenu from "./components/LeftMenu.tsx";
import { TopBanner } from "./components/TopBanner.tsx";
import { MainPanel } from "./components/MainPanel.tsx";

function App() {
    return (
        <PrimeReactProvider>
            <TopBanner></TopBanner>
            <div style={{ display: "flex",height: "calc(100vh - 5rem)"  }}>
                <LeftMenu></LeftMenu>
                <MainPanel></MainPanel>
            </div>

        </PrimeReactProvider>
    )
}

export default App