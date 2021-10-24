import { makeStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Link from '@material-ui/core/Link';
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';
import Box from '@material-ui/core/Box';


const useStyles = makeStyles((theme) => ({
  footerul: {
    margin: 0,
    padding: 0,
    listStyle: 'none',
  },
  appBar: {
    borderBottom: `1px solid ${theme.palette.divider}`,
  },
  toolbar: {
    flexWrap: 'wrap',
  },
  toolbarTitle: {
    flexGrow: 1,
  },
  link: {
    margin: theme.spacing(1, 1.5),
  },

  footer: {
    borderTop: `1px solid ${theme.palette.divider}`,
    marginTop: theme.spacing(8),
    paddingTop: theme.spacing(3),
    paddingBottom: theme.spacing(3),
    [theme.breakpoints.up('sm')]: {
      paddingTop: theme.spacing(6),
      paddingBottom: theme.spacing(6),
    },
  },
}));

// function Copyright() {
//   return (
//     <Typography variant="body2" color="textSecondary" align="center">
//       {'Copyright © '}
//       <Link color="inherit" href="https://material-ui.com/">
//         Johannes Heuel
//         </Link>{' '}
//       {new Date().getFullYear()}
//       {'.'}
//     </Typography>
//   );
// }

const footers = [
  {
    title: 'pollrBot',
    description: [
      {
        title: 'About',
        link: '/about'
      },
      {
        title: 'Tutorials',
        link: '/tutorials'
      },
    ]
  },
  {
    title: 'Legal',
    description: [
      {
        title: 'Legal Disclosure',
        link: '/legal_disclosure'
      },
      {
        title: 'Privacy policy',
        link: '/privacy'
      }, 
    ]
  }
];

function Footer() {
  const classes = useStyles();
  return (
    <Container maxWidth="md" component="footer" className={classes.footer}>
    <Grid container spacing={4} justify="space-evenly">
      {footers.map((footer) => (
        <Grid item xs={6} sm={3} key={footer.title}>
          <Typography variant="h6" color="textPrimary" gutterBottom>
            {footer.title}
          </Typography>
          <ul className={classes.footerul}>
            {footer.description.map((item) => (
              <li key={item.title}>
                <Link href={item.link} variant="subtitle1" color="textSecondary">
                  {item.title}
                </Link>
              </li>
            ))}
          </ul>
        </Grid>
      ))}
    </Grid>
    <Box mt={5}>
      {/* <Copyright /> */}
    </Box>
  </Container>
  );
}

export default Footer;


