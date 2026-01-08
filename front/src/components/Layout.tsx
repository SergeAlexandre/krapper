
import TopBar from './TopBar';
import Sidebar from './Sidebar';

const Layout = () => {
    return (
        <div className="flex flex-column h-screen">
            <TopBar />
            <div className="flex flex-grow-1 overflow-hidden">
                <div className="w-16rem flex-shrink-0 overflow-y-auto surface-0 border-right-1 surface-border">
                    <Sidebar />
                </div>
                <div className="flex-grow-1 overflow-y-auto p-4 surface-ground">
                    {/* Main content goes here */}
                    <h1>Welcome</h1>
                    <p>Select an item from the menu.</p>
                </div>
            </div>
        </div>
    );
}

export default Layout;
