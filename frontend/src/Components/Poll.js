import { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import Typography from '@material-ui/core/Typography';
import Paper from '@material-ui/core/Paper';
import { PieChart } from 'react-minimal-pie-chart';
import Grid from '@material-ui/core/Grid';
import Container from '@material-ui/core/Container';
import Button from '@material-ui/core/Button';
import React from 'react';


const colorScheme = [
  '#1f77b4',
  '#ff7f0e',
  '#2ca02c',
  '#d62728',
  '#9467bd',
  '#8c564b',
  '#e377c2',
  '#7f7f7f',
  '#bcbd22',
  '#17becf',
  '#9dc6b0',
  '#c6e0d2',
  '#b7a087',
  '#cfbeab',
  '#e0d4c6',
  '#524656',
  '#cf4647',
  '#eb7b59',
  '#e5ddcb',
  '#a7c5bd',
];

function percent(n, N) {
  return Math.round(100 * n / Math.max(1, N));
}

function Answer(props) {
  let answer = props.answer;

  if (answer == null) {
    answer = {
      'UserName': '...',
    }
  };
  if (props.flex == null) props.flex = false;

  if (props.flex) {
    return (
      <li key={props.idx} style={{
        listStyleType: 'none',
        listStylePosition: 'inside',
        wordWrap: 'break-word',
        maxWidth: '70vw',
        flex: 'auto',
//        border: '1px solid black',
        margin: '2px 5px',
      }}>
        <Button size='small' variant='outlined' style={{
          pointerEvents: 'none',
          userSelect: 'all',
          textTransform: 'none',
          display: 'inline-block',
          wordWrap: 'break-word',
          maxWidth: '80vw',
        }}>{answer['UserName']}</Button>
      </li>
    );
  }
}

function Answers(props) {
  const [showAll, setShowAll] = React.useState(false);

  if (props.answers == null) return null;
  if (props.answers.length === 0) return null;
  if (props.flex == null) props.flex = false;
  let n;
  let show_button_text;
  if (props.flex) {
    if (!showAll) {
      n = 15;
      show_button_text = 'Show all';
    } else {
      n = props.answers.length;
      show_button_text = 'Show less';
    }
    return (
      <div style={{
        paddingBottom: '2vh',
      }}>
          <ul style={{
          display: 'flex',
          flexWrap: 'wrap',
          justifyContent: 'space-between',
          listStyleType: 'none',
          padding: '3px',
          margin: '0',
        }}>
          {props.answers.slice(0, n).map((answer, index) => (
            <Answer answer={answer} key={index} idx={props.key + "_" + index} flex={props.flex} />
          ))}
        </ul>
        {props.answers.length > n &&
          <div style={{
            align: 'right',
            margin: '2px 5px',
          }}>
            <Button 
              size='small' 
              variant='text' 
              style={{
                // pointerEvents: 'none',
                userSelect: 'all',
                textTransform: 'none',
                display: 'inline-block',
                wordWrap: 'break-word',
                maxWidth: '80vw',
              }} 
              onClick={() => {setShowAll(!showAll)}}
            >+ {props.answers.length-n} more entries</Button>
          </div>
        }
        {(props.answers.length > n || showAll) &&
          <div style={{
            align: 'right',
            margin: '2px 5px',
          }}>
            <Button size='medium' variant='text' style={{
              display: 'block',
              marginLeft: 'auto', 
              marginRight: 'auto', 
            }} onClick={() => {setShowAll(!showAll)}}>{show_button_text}</Button>
          </div>
        }
      </div>);
  } else {
    return (
      <ul>
        {props.answers.map((answer, index) => (
          <Answer answer={answer} key={index} idx={props.key + "_" + index} />
        ))}
      </ul>
    );
  }
}

function Options(props) {
  if (props.options == null) return null;
  if (props.options.length === 0) return null;
  return (
    <div>
      {props.options.map((option, index) => (
        <div key={index}>
          <Typography variant="h6" style={{
            paddingBottom: "1vh",
          }}>{option['Text']} ({percent(option['Cnt'], props.total)}%)</Typography>
          <Answers answers={option['Answers']} idx={index} flex={true}/>
        </div>
      ))}
    </div>
  );
}

function Poll(props) {
  if (Object.keys(props.poll).length === 0) {
    return null;
  }
  return (
    <div style={{
      maxWidth: '600px',
    }}>
      <Typography variant="h5" style={{ paddingTop: '3vh', paddingBottom: '2vh' }}>
        📊 {props.poll['Question']}
      </Typography>

      <Options options={props.poll['Options']} total={props.poll['Cnt']} />
    </div>
  );
}

function PollPage() {
  let params = useParams();
  const [poll, setPoll] = useState({});

  const getPoll = () => {
    fetch('https://pollr.boosted.science/poll/' + params["id"]
      , {
        headers: {
          'Accept': 'application/json'
        }
      }
    )
      .then(function (response) {
        return response.json();
      })
      .then(function (myJson) {
        setPoll(myJson)
      });
  }

  useEffect(getPoll, [params]);


  return (
    <Paper elevation={0} style={{ 
        backgroundColor: '#fafafa', 
        marginTop: '20px',
      }}>
      <Container maxWidth="xl">
        <Grid 
          container 
          // spacing={2} 
          // justify="space-evenly"
          spacing={0}
          direction="column"
          alignItems="center"
          style={{ 
            minHeight: '100vh', 
            maxWidth: '100%'
          }}
          >
          <Grid item xs={12} md={12} xl={12}>
            <Chart poll={poll} />
          </Grid>
          <Grid item xs={12} md={12} xl={12}>
          {/* <Grid item xs={12} md={6} xl={4}> */}
            <Poll poll={poll} />
          </Grid>
          <Grid item xs={12} style={{marginTop: '40vh'}}>
            <a style={{display: 'block', width: '300', marginLeft: 'auto', marginRight: 'auto'}} href="https://s.click.aliexpress.com/e/_9xNflz?bz=300*250" target="_parent">
              <img width="300" height="250" alt="" src="https://ae01.alicdn.com/kf/HTB13jH6J4TpK1RjSZFKq6y2wXXaP/EN_300_250.jpg"/>
            </a>
          </Grid>
        </Grid>
      </Container>
    </Paper >
  );
}

function Chart(props) {
  if (Object.keys(props.poll).length === 0) {
    return null;
  }
  if (props.poll['Cnt'] < 2) {
    return null;
  }

  const data = [];
  props.poll['Options'].map((option, index) => (
    data.push({
      title: option['Text'],
      value: percent(option['Cnt'], props.poll['Cnt']),
      color: colorScheme[index]
    })
  ));
  const shiftSize = 7;

  return (
    <Grid item xs={12} md={12} xl={12}>
      <PieChart
        style={{
          display: 'block',
          marginLeft: 'auto',
          marginRight: 'auto',
          //maxWidth: '600',
          width: '70%'
        }}
        data={data}
        radius={PieChart.defaultProps.radius - shiftSize}
        segmentsShift={(index) => (2)}
        label={({ dataEntry }) => dataEntry.value > 0 ? dataEntry.value + '%' : ''}
        labelPosition={75}
        lineWidth={'100'}
        labelStyle={{
          fontSize: '7px',
          fontFamily: 'sans-serif'
        }}
      />
    </Grid>
  );
}
export default PollPage;



