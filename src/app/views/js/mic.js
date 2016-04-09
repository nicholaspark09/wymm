

function MicApp(){

    importScripts('/js/recordmp3.js');
    var that=this;
    var canvas;
    var audio;
    var modal;
    var Recorder;
    var audio_context;

    this.Init = function(){

      setupModal();


      try {
      // webkit shim
      window.AudioContext = window.AudioContext || window.webkitAudioContext;
      navigator.getUserMedia = ( navigator.getUserMedia ||
                       navigator.webkitGetUserMedia ||
                       navigator.mozGetUserMedia ||
                       navigator.msGetUserMedia);
      window.URL = window.URL || window.webkitURL;
      audio_context = new AudioContext;
      console.log('Audio context set up.');
      console.log('navigator.getUserMedia ' + (navigator.getUserMedia ? 'available.' : 'not present!'));
    } catch (e) {
      alert('No web audio support in this browser!');
    }
    navigator.getUserMedia({audio: true}, startUserMedia, function(e) {
      console.log('No live audio input: ' + e);
    });

//document.getElementById(canvasId);
      audio = $(document).find('#audio').get(0);//document.getElementById(videoId);
      $(document).on("click","#micButton", function(e){
        e.preventDefault();

      });

      $(document).on("click","#closeMic", function(e){
        e.preventDefault();
        recording = false;
      });

      $(document).on("click","#micsnap", function(e){
                // the fade only works on firefox?
                audio.src ="";
                worker.postMessage({command:'clear'});
                recording = true;
                $(this).text("Recording...");
                $(this).attr('class','btn btn-green');
            });

      $(document).on("click","#micreset", function(e){
                audio.src ="";
                worker.postMessage({command:'clear'});
            });

      $(document).on("click","#micstop", function(e){
                recording = false;
                that.stop();
                $('#micsnap').text("Record");
                $('#micsnap').attr('class','btn btn-danger');
                //var blob = createBlob();
            });
}


function startUserMedia(stream) {
    var input = audio_context.createMediaStreamSource(stream);
    console.log('Media stream created.' );
    console.log("input sample rate " +input.context.sampleRate);
    // Feedback!
    //input.connect(audio_context.destination);
    console.log('Input connected to audio context destination.');
    recorder = new Recorder(input, {
                  numChannels: 1
                });
    console.log('Recorder initialised.');
  }
/*

function interleave(leftChannel, rightChannel){
  var length = leftChannel.length + rightChannel.length;
  var result = new Float32Array(length);
 
  var inputIndex = 0;
 
  for (var index = 0; index < length; ){
    result[index++] = leftChannel[inputIndex];
    result[index++] = rightChannel[inputIndex];
    inputIndex++;
  }
  return result;
}



function writeUTFBytes(view, offset, string){ 
  var lng = string.length;
  for (var i = 0; i < lng; i++){
    view.setUint8(offset + i, string.charCodeAt(i));
  }
}
*/

    function bytesToSize(bytes) {
        var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
        if (bytes == 0) return 'n/a';
        var i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));
        if (i == 0) return bytes + ' ' + sizes[i];
          return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + sizes[i];
    };

    function setupModal(){

      var modal = '<div class="modal fade" id="audioModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">\
      <div class="modal-dialog" role="document">\
        <div class="modal-content">\
          <div class="modal-header">\
        <h4 class="modal-title" id="myModalLabel">Photo App</h4>\
      </div>\
      <div class="modal-body">\
            <audio src="" id="audio" controls autoplay></audio>\
      </div>\
      <div class="modal-footer">\
        <a href="#" type="button" class="btn btn-danger" id="micsnap">Record</a>\
        <a href="#" type="button" class="btn btn-danger" id="micstop">Stop</a>\
        <a href="#" type="button" class="btn btn-info" id="micreset">Reset</a>\
        <a href="#" type="button" class="btn btn-success" id="upload">Save</a>\
        <button type="button" class="btn btn-default" data-dismiss="modal" id="closeMic">Close</button>\
      </div>\
      </div>\
    </div>\
    </div>';
      $(modal).appendTo('body');
    }

    that.Init();
  }