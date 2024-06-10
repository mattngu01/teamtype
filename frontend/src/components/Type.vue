
<script lang="ts">
import { defineComponent } from 'vue';

// -1 = correct 
function firstDiffIndex(compare : string, truth : string) : number {
    for (let i = 0; i < compare.length; i++) {
        if (i > truth.length) {
            return truth.length
        }

        if (compare[i] != truth[i]) {
            return i
        }
    }

    return -1
}

export default defineComponent({
    data() {
        return {
            typedWord: "",
            quote: "We don't understand what really causes events to happen. History is the fiction we invent to persuade ourselves that events are knowable and that life has order and direction. That's why events are always reinterpreted when values change. We need new versions of history to allow for our current prejudices.",
            currentWordIndex: 0,
        }
    },
    computed: {
        quoteArray() : string[] {
            return this.quote.split(/(?<=\s)/);
        }
    },
    methods: {
        quoteHtml() : string {
            let correct : string = "";
            if (this.currentWordIndex > 0) {
                correct = "<span style=\"color: green\">" + this.quoteArray.slice(0, this.currentWordIndex).join("") + "</span>";
            }
            let htmlWord : string = "";
            let currentWord : string = this.quoteArray[this.currentWordIndex];

            let firstMistake : number = firstDiffIndex(this.typedWord, currentWord);

            if (firstMistake != -1) {
                htmlWord += "<span style=\"color: green\">" + currentWord.slice(0, firstMistake) + "</span>" 
                htmlWord += "<span style=\"background-color: red\">" + currentWord.slice(firstMistake, currentWord.length) + "</span>";
            } else {
                htmlWord += "<span style=\"color: green\">" + this.typedWord + "</span>";
                htmlWord += currentWord.slice(this.typedWord.length, currentWord.length);
            }


            let unfinished = this.quoteArray.slice(this.currentWordIndex+ 1, this.quoteArray.length).join("");

            return correct + htmlWord + unfinished
        }
    },
    watch: {
        typedWord(newTypedWord: string, oldTypedWord: string) {
            // based on assumption that this is called on each change
            let diff : number = Math.abs(oldTypedWord.length - newTypedWord.length);
            if (oldTypedWord.length < newTypedWord.length) {
                console.log("+" + newTypedWord.slice(-diff));
            } else {
                console.log("-" + oldTypedWord.slice(-diff));
            }

            if (newTypedWord == this.quoteArray[this.currentWordIndex])  {
                this.typedWord = "";
                this.currentWordIndex++;
                console.log("Finished word " + newTypedWord)
            }
        }
    }
})

</script>
<template>
    <div> Hello world </div>
    <div v-html="quoteHtml()"></div>
    <input for="test" spellcheck="false" autocapitalize="off" autocomplete="off" v-model="typedWord" placeholder="Type here">
</template>