import Typography from '@material-ui/core/Typography';
import Paper from '@material-ui/core/Paper';
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';
import { makeStyles } from '@material-ui/core/styles';
import logo from './logo.jpg';
import Button from '@material-ui/core/Button';
import Link from '@material-ui/core/Link';
import ChatBubble from './ChatBubble';



const useStyles = makeStyles((theme) => ({
  paper: {
    marginTop: theme.spacing(8),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  avatar: {
    margin: theme.spacing(1),
    borderRadius: '100%',
    width: '40vw',
    maxWidth: '150px'
  },
  heading: {
    paddingBottom: theme.spacing(2),
  },
}));

function Home() {
  const classes = useStyles();

  return (
    <Paper elevation={0} style={{ backgroundColor: '#fafafa' }}>
      <Container maxWidth="xl">
        <Grid container spacing={2} justify="space-evenly">
          <Grid item xs={12} md={6} xl={4} className={classes.paper}>
            <img className={classes.avatar} src={logo} alt="Logo" />
            <Typography component="h1" variant="h5" className={classes.heading}>
              pollrBot
            </Typography>
            <Typography style={{ width: '80%', maxWidth: 350, textAlign: 'center' }}>
              Hello! I am a chat bot for Telegram that helps you by creating polls for your group or private chat.
            </Typography>
            <Link style={{marginTop: '2vh', textDecoration: 'none'}} href="https://telegram.me/pollrBot">
              <Button variant="contained" color="primary">Add me on Telegram</Button>
            </Link>

            <Typography style={{marginTop: '50px'}}>
              Before creating polls, you need to activate the bot by sending "/start".
            </Typography>
            <ChatBubble
              messages={[
                {
                  type: 0,
                  text: '/start'
                },
                {
                  type: 1,
                  text: <>I can help you create, send and manage polls.<br/><br/>What do you want to do?</>,
                  buttons: ['create poll'],
                },
              ]}
            />
            <Typography style={{marginTop: '20px'}}>
              After a press on the "create poll" button the poll creation is started, asking you to first give a question and then a number of options.
            </Typography>
            <ChatBubble
              messages={[
                {
                  type: 1,
                  text: <>Great! Send a question for the new poll, please.</>
                },
                {
                  type: 0,
                  text: 'When do we meet?'
                },
                {
                  type: 1,
                  text: <>OK now that we got a question, please send answer options to your poll.</>
                },
                {
                  type: 0,
                  text: 'today'
                },
                {
                  type: 0,
                  text: 'tomorrow'
                },
                {
                  type: 1,
                  text: <>
                      You can add more options by sending messages each containing one option. If you are done, please push the done button.
                      <br/><br/>
                      Preview:<br/>
                      When do we meet?<br/>
                      1) today<br/>
                      2) tomorrow
                    </>,
                  buttons: ['done'],
                },
              ]}
            />
            <Typography style={{marginTop: '20px'}}>
              When all options are added, the poll can be finalized by pressing the "done" button. The bot will send a final preview of the poll with a share button attached to the message.
            </Typography>
            <ChatBubble
              messages={[
                {
                  type: 1,
                  text: <>
                      Finished creating a new poll.
                      <br/><br/>
                      Preview:<br/>
                      When do we meet?<br/>
                      1) today<br/>
                      2) tomorrow
                    </>,
                  buttons: ['Share poll', 'create new poll'],
                },
              ]}
            />
            <Typography style={{marginTop: '20px'}}>
              The shared poll message in a group or private chat will then include a poll listing with percentages of how many people voted for the options, a list of users under every poll option, a counter how many people voted and a link to a web view of the poll.
            </Typography>
            <ChatBubble
              messages={[{
                type: 0,
                text: <>
                    📊 When do we meet?
                    <br/><br/>
                    today (100%):<br/>
                    └ Johannes Heuel<br/><br/>
                    tomorrow (100%):<br/><br/>
                    1 👥  <a href="https://pollrbot.com/p/G/L1jGgb">pollrbot.com/p/‎G/L1jGgb</a>
                  </>,
                buttons: ['today (1)', 'tomorrow (0)'],
              },
            ]}
            />
            <Typography style={{marginTop: '20px'}}>
              Also these <Link href="/tutorials">video tutorials</Link> might help.
            </Typography>
          </Grid>
        </Grid>
      </Container>
    </Paper>
  )
}


export default Home;