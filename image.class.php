<?php
	class Image {
		private $img = array();
		public $oldimg = array(); 
		public $newimg = array();  
		
		function Image($imgName) {
			$this->img['name'] = basename($imgName); 
			$a = explode(".", $this->img['name']); 
			$this->img['short_name'] = $a[0]; 
			$this->img['filepath'] = $imgName; 
			$this->img['size'] = getimagesize($imgName);  
		}

		function RenameImage($old, $new) {
			system("mv " . $old. " " . $new); 
		}

		function GetShortName() {
			return $this->img['short_name']; 
		}

		function GetSize() {
			return $this->img['size']; 
		}
		
		function GetName() {
			return $this->img['name']; 
		}

		function GetNameNoExt() {
			$p = pathinfo($this->img['name']); 
			return $p['filename']; 
		}

		function GetPath() {
			return $this->img['filepath']; 
		}

		function ScaleImage($dest, $factor) {
			// $ratio = ($this->img['size'][0] / $this->img['size'][1]) * ($factor / 100); 
			$newWidth = $this->img['size'][0] * ($factor / 100); 
			$newHeight = $this->img['size'][1] * ($factor / 100); 
			
			$this->ResizeImage($dest, $newWidth, $newHeight); 
		}	
		
		function ResizeImage($dest, $newWidth, $newHeight) {
			$size = $this->img['size']; 
			$image_template = imagecreatetruecolor($newWidth, $newHeight);

			$imgType = array_reverse(explode(".", $this->img['filepath'])); 
			switch ($imgType[0]) { 
				case "jpg": 
					$imgSrc = imagecreatefromjpeg($this->img['filepath']) or die("Is " . $this->img['name'] . " really a jpg image?\n"); 
				case "png": 
					$imgSrc = imagecreatefrompng($this->img['filepath']) or die("Is " . $this->img['name'] . " really a png image?\n");
			} 
			
			$newImg = imagecopyresampled($image_template, $imgSrc, 0, 0, 0, 0, $newWidth, $newHeight, $size[0], $size[1]); 
			
			$result = imagejpeg($image_template, $dest, 100); 
			
			imagedestroy($image_template); 

			$this->oldimg['height'] = $size[1]; 
			$this->oldimg['width']  = $size[0]; 
			$this->newimg['height'] = $newHeight; 
			$this->newimg['width']  = $newWidth; 
		}

		function CopyImage($dest) {
			$image_template = imagecreatetruecolor($this->img['size'][0], $this->img['size'][1]); 
			$imgSrc = imagecreatefromjpeg($this->img['filepath']); 
			$newImg = imagecopyresampled($image_template, $imgSrc, 0, 0, 0, 0, $this->img['size'][0], $this->img['size'][1], $this->img['size'][0], $this->img['size'][1]); 
			$result = imagejpeg($image_template, $dest, 100); 

			imagedestroy($image_template); 
		}

		function MakeSmall($dest) {
			$this->ResizeImage($dest, 200, 160); 
		}

		function MakeTN($dest) {
			$this->ResizeImage($dest, 110, 90); 
		}

		function MakeReg($dest) {
			$this->ResizeImage($dest, 380, 285); 
		}
	}
?>
