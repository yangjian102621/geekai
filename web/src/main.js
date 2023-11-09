import {createApp} from 'vue'
import ElementPlus from "element-plus"
import "element-plus/dist/index.css"
import 'vant/lib/index.css';
import App from './App.vue'
import {createPinia} from "pinia";
import {
    Button,
    Cell,
    CellGroup,
    ConfigProvider,
    Dialog,
    DropdownItem,
    DropdownMenu,
    Field,
    Form,
    Icon,
    Image,
    List,
    NavBar,
    Notify,
    Picker,
    Popup,
    Search,
    ShareSheet,
    Sticky,
    SwipeCell,
    Switch,
    Tabbar,
    TabbarItem,
    Tag,
    TextEllipsis,
    Uploader
} from "vant";
import {router} from "@/router";
import 'v3-waterfall/dist/style.css'
import V3waterfall from "v3-waterfall";

const app = createApp(App)
app.use(createPinia())
app.use(ConfigProvider);
app.use(Tabbar);
app.use(TabbarItem);
app.use(NavBar);
app.use(Search);
app.use(Cell)
app.use(Image)
app.use(TextEllipsis)
app.use(Notify)
app.use(Picker)
app.use(Popup)
app.use(List);
app.use(Form);
app.use(Field);
app.use(CellGroup);
app.use(Button);
app.use(DropdownMenu);
app.use(Icon);
app.use(DropdownItem);
app.use(Sticky);
app.use(SwipeCell);
app.use(Dialog);
app.use(ShareSheet);
app.use(Switch);
app.use(Uploader);
app.use(Tag);
app.use(V3waterfall)
app.use(router).use(ElementPlus).mount('#app')


