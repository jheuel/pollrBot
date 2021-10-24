import {
  Router,
  Switch,
  Route,
} from 'react-router-dom';
import { useEffect } from 'react';

import { createBrowserHistory } from 'history';
import ReactGA from 'react-ga';
import PollPage from './Components/Poll';
import Home from './Components/Home';
import Tutorials from './Components/Tutorials';
import LegalDisclosure from './Components/LegalDisclosure';
import PrivacyPolicy from './Components/PrivacyPolicy';
import Footer from './Components/Footer.js';
import About from './Components/About.js';
import Contact from './Components/Contact.js';
import NotFound from './Components/NotFound.js';

import React from 'react';
import AppBar from '@material-ui/core/AppBar';
import CssBaseline from '@material-ui/core/CssBaseline';
import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
// import Hidden from '@material-ui/core/Hidden';
import IconButton from '@material-ui/core/IconButton';
import InboxIcon from '@material-ui/icons/MoveToInbox';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import MailIcon from '@material-ui/icons/Mail';
import MenuIcon from '@material-ui/icons/Menu';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import { createMuiTheme, makeStyles, ThemeProvider } from '@material-ui/core/styles';
import Grid from '@material-ui/core/Grid';
import { Link as RouterLink } from 'react-router-dom';
import Link from '@material-ui/core/Link';



const drawerWidth = 240;

const useStyles = makeStyles((theme) => ({
  root: {
    display: 'flex',
  },
  // drawer: {
  //   [theme.breakpoints.up('md')]: {
  //     width: drawerWidth,
  //     flexShrink: 0,
  //   },
  // },
  // appBar: {
  //   [theme.breakpoints.up('md')]: {
  //     width: `calc(100% - ${drawerWidth}px)`,
  //     marginLeft: drawerWidth,
  //   },
  // },
  // menuButton: {
  //   marginRight: theme.spacing(2),
  //   [theme.breakpoints.up('md')]: {
  //     display: 'none',
  //   },
  // },
  // necessary for content to be below app bar
  toolbar: theme.mixins.toolbar,
  drawerPaper: {
    width: drawerWidth,
  },
  content: {
    flexGrow: 1,
    padding: theme.spacing(2),
    paddingTop: theme.spacing(2),
  },
}));

ReactGA.initialize('UA-187821293-1');

const history = createBrowserHistory();
history.listen((location) => {
  ReactGA.pageview(location.pathname + location.search);
});

const theme = createMuiTheme({

});

function App(props) {
  const { window } = props;
  const classes = useStyles();
  //const theme = useTheme();
  const [mobileOpen, setMobileOpen] = React.useState(false);

  useEffect(() => {
    if (window !== undefined) {
      ReactGA.pageview(window.location.pathname + window.location.search);
    }
  })

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const navItems = [
    {
      title: 'Home',
      url: '/'
    },
    {
      title: 'About',
      url: '/about'
    },
    {
      title: 'Tutorials',
      url: '/tutorials'
    },
    {
      title: 'Contact',
      url: '/contact'
    }
  ];
  const drawer = (
    <div>
      <div className={classes.toolbar} />
      <Divider />
      <List>
        {navItems.map((item, index) => (
          <ListItem button component={RouterLink} to={item.url} key={item.title} onClick={handleDrawerToggle}>
            <ListItemIcon>{index % 2 === 0 ? <InboxIcon /> : <MailIcon />}</ListItemIcon>
            <ListItemText primary={item.title} />
          </ListItem>
        ))}
        <Divider />
      </List>

    </div>
  );

  const container = window !== undefined ? () => window().document.body : undefined;

  return (
    <Router history={history}>
      <div className={classes.root}>
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <AppBar position="fixed" className={classes.appBar}>
            <Toolbar>
              <IconButton
                color="inherit"
                aria-label="open drawer"
                edge="start"
                onClick={handleDrawerToggle}
                className={classes.menuButton}
              >
                <MenuIcon />
              </IconButton>
              <Link href="https://telegram.me/pollrBot" color="inherit" style={{ textDecoration: 'none' }}>
                <Typography variant="h6" noWrap>
                  pollrBot
          </Typography>
              </Link>
            </Toolbar>
          </AppBar>
          <nav className={classes.drawer} aria-label="mailbox folders">
            {/* The implementation can be swapped with js to avoid SEO duplication of links. */}
            {/* <Hidden mdUp implementation="js"> */}
            <Drawer
              container={container}
              variant="temporary"
              anchor={theme.direction === 'rtl' ? 'right' : 'left'}
              open={mobileOpen}
              onClose={handleDrawerToggle}
              classes={{
                paper: classes.drawerPaper,
              }}
              ModalProps={{
                keepMounted: true, // Better open performance on mobile.
              }}
            >
              {drawer}
            </Drawer>
            {/* </Hidden> */}
            {/* <Hidden smDown implementation="js">
            <Drawer
              classes={{
                paper: classes.drawerPaper,
              }}
              variant="permanent"
              open
            >
              {drawer}
            </Drawer>
          </Hidden> */}
          </nav>
          <main className={classes.content}>
            <div className={classes.toolbar} />
            <Grid container spacing={0} style={{ minHeight: '85vh' }}>
              <Grid item xs={12}>
                <Switch>
                  <Route exact path="/" component={Home} />
                  <Route path="/p/:id+" component={PollPage} />
                  <Route path="/privacy" component={PrivacyPolicy} />
                  <Route path="/legal_disclosure" component={LegalDisclosure} />
                  <Route path="/tutorials" component={Tutorials} />
                  <Route path="/about" component={About} />
                  <Route path="/contact" component={Contact} />
                  <Route component={NotFound} />
                </Switch>
              </Grid>
            </Grid>
            <Footer />
          </main>
        </ThemeProvider>
      </div>
    </Router>
  );
}

export default App;
