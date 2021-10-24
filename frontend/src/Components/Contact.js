import Typography from '@material-ui/core/Typography';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Grid from '@material-ui/core/Grid';
import MailOutlineIcon from '@material-ui/icons/MailOutline';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';


const useStyles = makeStyles((theme) => ({
  paper: {
    marginTop: theme.spacing(8),
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
  },
  avatar: {
    margin: theme.spacing(1),
    backgroundColor: theme.palette.secondary.main,
  },
  form: {
    width: '100%', // Fix IE 11 issue.
    marginTop: theme.spacing(3),
  },
  submit: {
    margin: theme.spacing(3, 0, 2),
  },
}));


function Contact() {
  const classes = useStyles();
  const [firstName, setFirstName] = React.useState('');
  const [lastName, setLastName] = React.useState('');
  const [email, setEmail] = React.useState('');
  const [message, setMessage] = React.useState('');
  const [send, setSend] = React.useState(false);

  const submitForm = (event) => {
    event.preventDefault();
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        FirstName: firstName,
        LastName: lastName,
        Email: email,
        Message: message,
      })
    };
    fetch('https://pollrbot.com/messageme/', requestOptions)
      .then((response) => response.json())
      .then((data) => setSend(data.Success));
  };

  return (
    <Paper elevation={0} style={{ backgroundColor: '#fafafa' }}>

      <Container component="main" maxWidth="xs">
        <CssBaseline />
        {!send &&
          <div className={classes.paper}>
            <Avatar className={classes.avatar}>
              <MailOutlineIcon />
            </Avatar>
            <Typography component="h1" variant="h5">
              Contact
            </Typography>
            <form className={classes.form} noValidate onSubmit={(event) => (event.preventDefault())}>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6}>
                  <TextField
                    autoComplete="fname"
                    name="firstName"
                    variant="outlined"
                    required
                    fullWidth
                    id="firstName"
                    label="First Name"
                    autoFocus
                    onSubmit
                    onChange={(e) => (setFirstName(e.target.value))}
                  />
                </Grid>
                <Grid item xs={12} sm={6}>
                  <TextField
                    variant="outlined"
                    required
                    fullWidth
                    id="lastName"
                    label="Last Name"
                    name="lastName"
                    autoComplete="lname"
                    onChange={(e) => (setLastName(e.target.value))}
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    variant="outlined"
                    required
                    fullWidth
                    id="email"
                    label="Email Address"
                    name="email"
                    autoComplete="email"
                    onChange={(e) => (setEmail(e.target.value))}
                  />
                </Grid>
                <Grid item xs={12}>
                  <TextField
                    variant="outlined"
                    required
                    fullWidth
                    name="message"
                    label="Message"
                    type="message"
                    id="message"
                    multiline
                    rows={4}
                    rowsMax={20}
                    onChange={(e) => (setMessage(e.target.value))}
                  />
                </Grid>
                <Grid item xs={12}>
                </Grid>
              </Grid>
              <Button
                fullWidth
                variant="contained"
                color="primary"
                className={classes.submit}
                onClick={submitForm}
              >
                Send message
              </Button>
            </form>
          </div>
        }
        {send &&
          <div className={classes.paper}>
            <Typography variant="h4">Thank you!</Typography>

            <Typography>Your message has been sent.</Typography>
          </div>
        }
      </Container>
    </Paper>
  )
}


export default Contact;