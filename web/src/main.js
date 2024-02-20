import {createApp} from 'vue'
import ElementPlus from "element-plus"
import "element-plus/dist/index.css"
import 'vant/lib/index.css';
import App from './App.vue'
import {createPinia} from "pinia";
import {
    Badge,
    Button,
    Cell,
    CellGroup,
    Circle,
    Col,
    Collapse,
    CollapseItem,
    ConfigProvider,
    Dialog,
    DropdownItem,
    DropdownMenu,
    Empty,
    Field,
    Form,
    Grid,
    GridItem,
    Icon,
    Image,
    ImagePreview,
    Lazyload,
    List,
    Loading,
    NavBar,
    Notify,
    Overlay,
    Picker,
    Popup,
    Row,
    Search,
    ShareSheet,
    Slider,
    Sticky,
    SwipeCell,
    Switch,
    Tab,
    Tabbar,
    TabbarItem,
    Tabs,
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
app.use(Overlay)
app.use(Col)
app.use(Row)
app.use(Slider)
app.use(Badge)
app.use(Collapse);
app.use(CollapseItem);
app.use(Grid);
app.use(GridItem);
app.use(Empty);
app.use(Circle);
app.use(Loading);
app.use(Lazyload);
app.use(ImagePreview);
app.use(Tab);
app.use(Tabs);
app.use(router).use(ElementPlus).mount('#app')


