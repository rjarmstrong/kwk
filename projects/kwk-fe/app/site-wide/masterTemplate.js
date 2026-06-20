import {Link} from 'inferno-router'
import {Logo} from './svg';
import {decodeToken} from '../api/client'
import FindBox from './findBox'
import Creator from './creator'
import feather from 'feather-icons'
import ProfileBadge from './profileBadge'
import SnpContext from './snpContext'
import {info} from "./util";
import str from "../stores/snpStore";

window.feather = feather;

export default class MasterTemplate extends Component {

    constructor({props, params}) {
        super(props);
        this.params = params || {};
        let token = decodeToken() || {};
        this.home = '/';
        if (token.key !== 'anon') {
            this.home += token.key;
        }
    }

    componentDidMount() {
        document.addEventListener('contextmenu', function (e) {
            info('close existing');
            str.notify('invoke-snp-context', 'snpContext', { visible : false });
            e.preventDefault();
        }, true);
        document.addEventListener('click', function (e) {
            str.notify('invoke-snp-context', 'snpContext', { visible : false });
            e.preventDefault();
        }, false);
    }

    render({ children }) {
        return (<div className="struct">
                <div className="navbar">
                    <div className="navbar-start">
                        <Link to={this.home}>
                            <Logo/>
                        </Link>
                        <div className="navbar-item">
                            <FindBox u={this.params.u}/>
                        </div>
                    </div>
                    <div className="navbar-end">
                        <div className="navbar-item">
                            <Creator/>
                        </div>
                        <div className="navbar-item">
                            <ProfileBadge/>
                        </div>
                    </div>
                </div>
                <div className="section main struct">
                    {children}
                    <script>
                        feather.replace();
                    </script>
                </div>
                <SnpContext/>
            </div>
        )
    }
}
/*
| Timeline | Themes | Extensions | Usage
*/