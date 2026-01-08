import './App.css'
import {PrimeReactProvider} from 'primereact/api';
import {Button} from "primereact/button";

import "primereact/resources/themes/lara-light-cyan/theme.css";
//import 'primereact/resources/themes/mdc-dark-indigo/theme.css';

function App() {
    return (
        <PrimeReactProvider>
            <Button>Hello</Button>
        </PrimeReactProvider>
    )
}

export default App