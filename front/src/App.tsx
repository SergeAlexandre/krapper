import './App.css'
import {PrimeReactProvider} from 'primereact/api';


import 'primeicons/primeicons.css';

import "primereact/resources/themes/lara-light-cyan/theme.css";
import LeftMenu from "./components/LeftMenu.tsx";
import {TopBanner} from "./components/TopBanner.tsx";
import {MainPanel} from "./components/MainPanel.tsx";
//import 'primereact/resources/themes/mdc-dark-indigo/theme.css';

function App() {
    return (
        <PrimeReactProvider>
            <TopBanner></TopBanner>
            <LeftMenu></LeftMenu>
            <MainPanel></MainPanel>

        </PrimeReactProvider>
    )
}

export default App