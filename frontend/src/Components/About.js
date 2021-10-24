import Typography from '@material-ui/core/Typography';
import Paper from '@material-ui/core/Paper';
import Link from '@material-ui/core/Link';
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';
import { makeStyles } from '@material-ui/core/styles';


const useStyles = makeStyles((theme) => ({
  paper: {
    marginTop: theme.spacing(8),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  heading: {
    paddingBottom: theme.spacing(2),
  },
}));

function About() {
  const classes = useStyles();


  return (
    <Paper elevation={0} style={{ backgroundColor: '#fafafa' }}>
      <Container maxWidth="xl">
        <Grid container spacing={2} justify="space-evenly">
          <Grid item xs={12} md={6} xl={4} className={classes.paper}>
            <Typography component="h1" variant="h5" className={classes.heading}>
                About
            </Typography>
            <Typography style={{textAlign: 'left'}}>
              This is the website for pollrBot, a Telegram bot that is great at creating polls for group and normal chats. You can add pollrBot <Link href="https://telegram.me/pollrBot">here</Link>.
              If you need help to get started these <Link href="/tutorials">video tutorials</Link> might help.
            </Typography>
            <Typography style={{textAlign: 'left', marginTop: '16px'}}>
              This bot started as a way to decide on a pizza place where to order food.
              However, as it seems that it works quite well, more and more users started using it (now more than a million users total).
              It is used for simple polls, multiple-choice polls and for example as a way of documenting attendance.
            </Typography>
          </Grid>
        </Grid>
      </Container>
    </Paper>
  )
}

export default About;