import Paper from '@material-ui/core/Paper';
import Typography from '@material-ui/core/Typography';
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';
import React from 'react';
import Card from '@material-ui/core/Card';
import CardActionArea from '@material-ui/core/CardActionArea';
import CardActions from '@material-ui/core/CardActions';
import CardContent from '@material-ui/core/CardContent';
import Button from '@material-ui/core/Button';
import Link from '@material-ui/core/Link';


function Tutorials() {
  const videos = [
    {
      title: "Mr. Hengki Debora",
      url: "https://www.youtube.com/embed/OMvaM0i7ZH8",
    },
    {
      title: "Steve Berthold",
      url: "https://www.youtube.com/embed/2VArGXVYOg4",
    },
    {
      title: "Yulianti_Zu Channel",
      url: "https://www.youtube.com/embed/MY2XybLtYc0",
    },
    {
      title: "Ibnu Fajar75 CHANNEL",
      url: "https://www.youtube.com/embed/c7oIyrlAqKM",
    },
    {
      title: "Andi Subhan",
      url: "https://www.youtube.com/embed/vfSomqWT3EQ",
    },
    {
      title: "R. A. Selvi Saptya Dewi",
      url: "https://www.youtube.com/embed/bveZt7WXkNY",
    },
    {
      title: "Adien Novarisa",
      url: "https://www.youtube.com/embed/8RLmKZaF__M",
    },
    {
      title: "kuntilanak jongang",
      url: "https://www.youtube.com/embed/nP_AqF7CJwk",
    },
    {
      title: "TheFundamentalist",
      url: "https://www.youtube.com/embed/k6DFDcrTL64",
    },
    {
      title: "pollrBot tutrial Aan Andriana",
      url: "https://www.youtube.com/embed/41u1fGmj0oc",
    },
    {
      title: "ApanatorHD",
      url: "https://www.youtube.com/embed/iBLflyRk4TQ",
    },
    {
      title: "Rino Sugiantoro",
      url: "https://www.youtube.com/embed/FNKMfazSoIg",
    },
    {
      title: "aiGez TUTORIAL",
      url: "https://www.youtube.com/embed/idehwQm4_84",
    },
  ]

  const width = 560;
  const height = 315;

  return (
    <Paper elevation={0} style={{ backgroundColor: '#fafafa' }}>
      <Typography>Thanks to a few nice YouTubers there are video tutorials that show in detail how you can create your first poll</Typography>
      <Container maxWidth="xl">
        <Grid container spacing={2} justify="space-evenly">
          {videos.map((video) => (
            <Grid item xs={12} md={6} xl={4} key={video.title}>
              <Card>
                <CardActionArea>
                  {video.frame}
                  <CardContent>
                    <Typography gutterBottom variant="h5" component="h2">
                      <div style={{
                        position: 'relative',
                        paddingBottom: '56.25%',
                        paddingTop: '25px',
                        height: '0',
                      }}>
                        <iframe title={video.title} style={{
                          position: 'absolute',
                          top: '0',
                          left: '0',
                          width: '100%',
                          height: '100%',
                        }} frameborder="0" allowfullscreen width={width} height={height} src={video.url}></iframe>
                      </div>
                    </Typography>
                  </CardContent>
                </CardActionArea>
                <CardActions>
                  <Link href={video.url}>
                    <Button size="small" color="primary">
                      Watch on YouTube
                    </Button>
                  </Link>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Container>

    </Paper>
  )
}

export default Tutorials;