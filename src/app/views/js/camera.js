function CameraApp(){
    var that=this;
    var canvas;
    var video;
    var ctx;
    var image_format;
    var jpeg_quality;
    var errBack;
    var videoObj;
    var localStream;

    this.Init = function(){

      setupModal();


      canvas = $(document).find('#my_result').get(0);//document.getElementById(canvasId);
      video = $(document).find('#video').get(0);//document.getElementById(videoId);
      ctx = canvas.getContext("2d");
      image_format= "jpeg";
      jpeg_quality= 85;
      videoObj = { "video": true };
                errBack = function(error) {
                    console.log("Video capture error: ", error.code); 
                };

      $(document).on("click","#cameraButton", function(e){
        e.preventDefault();
        if(navigator.getUserMedia) { // Standard
                navigator.getUserMedia(videoObj, function(stream) {
                    video.src = stream;
                    video.play();
                    localStream = stream;
                    $("#snap").show();
                }, errBack);
            } else if(navigator.webkitGetUserMedia) { // WebKit-prefixed
                navigator.webkitGetUserMedia(videoObj, function(stream){
                    video.src = window.webkitURL.createObjectURL(stream);
                    video.play();
                    localStream = stream;
                    $("#snap").show();
                }, errBack);
            } else if(navigator.mozGetUserMedia) { // moz-prefixed
                navigator.mozGetUserMedia(videoObj, function(stream){
                    video.src = window.URL.createObjectURL(stream);
                    video.play();
                    localStream = stream;
                    $("#snap").show();
                }, errBack);
            }
      });

      $(document).on("click","#closeCamera", function(e){
        e.preventDefault();
        localStream.stop(); 
      });

      $(document).on("click","#snap", function(e){
                ctx.drawImage(video, 0, 0, 320, 240);
                // the fade only works on firefox?
                $("#video").fadeOut("slow");
                $("#my_result").fadeIn("slow");
                $("#snap").hide();
                $("#reset").show();
            });

      $(document).on("click","#reset", function(e){
                $("#video").fadeIn("slow");
                $("#my_result").fadeOut("slow");
                $("#snap").show();
                $("#reset").hide();
            });
      /*
      document.getElementById("upload").addEventListener("click", function(){
                var dataUrl = canvas.toDataURL("image/jpeg", 0.85);
               var bytes = encodeURI(dataUrl).split(/%..|./).length - 1;
               console.log(dataUrl);
               $.ajax({
                  type: 'POST',
                  url: '/multimedia/profilepic',
                  data:{
                    'controller':'profiles',
                    'safekey':safekey,
                    'name':"Picture",
                    'description':'',
                    'body':dataUrl,
                    'type':1
                  },
                  success: function(data){
                    var results= $.parseJSON(data);
                    if(results['Result']=="Success"){
                      localStream.stop();
                      alert("Saved!");
                    }else
                      alert(results['Error']);
                  },
                  error: function(){
                    console.log("No connection");
                  }
               });
      }); 

    */
 
    }

    this.getData = function(){
      var dataUrl = canvas.toDataURL("image/jpeg", 0.85);
      return dataUrl;
    }



    function bytesToSize(bytes) {
        var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
        if (bytes == 0) return 'n/a';
        var i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));
        if (i == 0) return bytes + ' ' + sizes[i];
          return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + sizes[i];
    };

    function setupModal(){

      var modal = '<div class="modal fade" id="photoModal" tabindex="-1" role="dialog" aria-labelledby="myModalLabel">\
      <div class="modal-dialog" role="document">\
        <div class="modal-content">\
          <div class="modal-header">\
        <h4 class="modal-title" id="myModalLabel">Photo App</h4>\
      </div>\
      <div class="modal-body">\
            <video id="video" autoplay width="320" height="240"></video>\
            <canvas width="320" height="240" id="my_result" style="display:none;"></canvas>\
      </div>\
      <div class="modal-footer">\
        <a href="#" type="button" class="btn btn-danger" id="snap">Take Photo</a>\
        <a href="#" type="button" class="btn btn-info" id="reset">Reset</a>\
        <a href="#" type="button" class="btn btn-success" id="upload">Save</a>\
        <button type="button" class="btn btn-default" data-dismiss="modal" id="closeCamera">Close</button>\
      </div>\
      </div>\
    </div>\
    </div>';
      $(modal).appendTo('body');
    }

    that.Init();
  }